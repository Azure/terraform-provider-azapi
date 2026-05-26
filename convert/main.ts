import path from 'path';
import { appendFileSync, readFileSync, readdirSync, writeFileSync } from 'fs';
import * as hcl from "hcl2-parser"
import { Bicep } from 'bicep-node';
import { series } from 'async';

const examples = path.join(__dirname, '../examples');

function findAllMainTfFiles() {
  let results: string[] = [];
  for (const entry of readdirSync(examples, { withFileTypes: true, recursive: true })) {
    if (entry.isFile() && 
      entry.name === 'main.tf' &&
      /[^\//]+@[^\//]+$/.test(entry.parentPath)) {
      results.push(path.join(entry.parentPath, entry.name));
    }
  }
  return results;
}

function tfToBicep(input: string): string {
  let output = '';
  const parsed = hcl.parseToObject(input);
  if (!parsed[0]) {
    throw new Error('Invalid HCL input or no resources found.');
  }

  const variables: Record<string, any> = parsed[0].variable;
  const resources: Record<string, any> = parsed[0].resource.azapi_resource;

  const indentLevel = 0;
  if (variables) {
    for (const [name, value] of Object.entries(variables)) {
      const defaultValue = convertToBicep(value[0].default, indentLevel);
      const type = value[0].type === '${string}' ? 'string' : null;

      output += line(`param ${name} ${type} = ${defaultValue}`);
    }
    output += line('');
  }

  const validTopLevelProperties = ['name', 'location', 'tags', 'sku', 'kind', 'managedBy', 'managedByExtended', 'extendedLocation', 'zones', 'plan', 'scale', 'identity'];

  for (const [name, value] of Object.entries(resources)) {
    const type = value[0].type;
    if (type.startsWith('Microsoft.Resources/resourceGroups@')) {
      continue;
    }

    output += line(`resource ${name} '${type}' = {`);

    const parentMatch = value[0].parent_id.match(/^\$\{azapi_resource\.([^.]+)\.id\}$/);
    if (parentMatch && parentMatch[1] !== 'resourceGroup') {
      output += line(indent(indentLevel + 2, `parent: ${parentMatch[1]}`));
    }

    const topLevelProps = Object.entries(value[0]).filter(x => validTopLevelProperties.includes(x[0]));
    for (const [propName, propValue] of topLevelProps) {
      output += line(indent(indentLevel + 2, objectProperty(propName, propValue, indentLevel + 2)));
    }
 
    const bodyProps = Object.entries(value[0].body ?? {});
    for (const [bodyPropName, bodyPropValue] of bodyProps) {
      output += line(indent(indentLevel + 2, objectProperty(bodyPropName, bodyPropValue, indentLevel + 2)));
    }

    output += line(line(`}`));
  }


  return output;
}

function convertToBicep(input: any, indentLevel: number): string {
  if (input === null || input === undefined) {
    return 'null';
  }

  if (Array.isArray(input)) {
    if (input.length === 0) {
      return '[]';
    }

    return [
      '[',
      ...input.map(item => indent(indentLevel + 2, convertToBicep(item, indentLevel + 2))),
      indent(indentLevel, ']')
    ].join('\n');
  }

  if (typeof input === 'object') {
    if (Object.keys(input).length === 0) {
      return '{}';
    }

    return [
      '{',
      ...Object.entries(input).map(
        ([key, value]) => indent(indentLevel + 2, objectProperty(key, value, indentLevel + 2))
      ),
      indent(indentLevel, '}')
    ].join('\n');
  }

  if (typeof input === 'number') {
    return Number.isInteger(input) ? `${input}` : `json('${input}')`;
  }

  if (typeof input === 'boolean') {
    return `${input}`;
  }

  if (typeof input === 'string') {
    let str = input
      .replace(/\${var\.([a-zA-Z0-9_]+)}/g, '${$1}')
      .replace(/\${azapi_resource\.(\w+)\.output(\.\w+(\.\w+)+)}/g, '${$1$2}')
      .replace(/\${azapi_resource\.(\w+(\.\w+)+)}/g, '${$1}')
      .replace(/\${resourceGroup\.location}/g, '${resourceGroup().location}');

    if (str.startsWith('${') && str.endsWith('}')) {
      return str.slice(2, -1);
    }

    str = str
      .replace(/\\/g, '\\\\')
      .replace(/\r/g, '\\r')
      .replace(/\n/g, '\\n')
      .replace(/\t/g, '\\t')
      .replace(/'/g, "\\'");
    return `'${str}'`;
  }

  throw new Error(`Unsupported type: ${typeof input} for value: ${input}`);
}

function indent(level: number, input: string) {
  return ' '.repeat(level) + input;
}

function line(input: string) {
  return input + '\n';
}

function objectProperty(key: string, value: any, indentLevel: number): string {
  const escaped = /^[a-zA-Z_][a-zA-Z0-9_]*$/.test(key) ? key : `'${key}'`;
  return `${escaped}: ${convertToBicep(value, indentLevel)}`;
}

async function main() {
  let markdown = '# Conversion Results\n';
  const bicepCliPath = await Bicep.install(path.join(__dirname, 'tmp'));
  const bicep = await Bicep.initialize(bicepCliPath);

  for (const tfPath of findAllMainTfFiles()) {
    const bicepPath = tfPath.replace(/\.tf$/, '.bicep');
    const relBicepPath = path.relative(examples, bicepPath).replace(/\\/g, '/');
    markdown += `## [${relBicepPath}](../examples/${relBicepPath})\n`;
    
    try {
      const tfContents = readFileSync(tfPath, 'utf8');
      const bicepContents = tfToBicep(tfContents);

      writeFileSync(bicepPath, bicepContents, 'utf8');

      const compileResult = await bicep.compile({ path: bicepPath });
      const diagnostics = compileResult.diagnostics
        .map(d => `[${d.level} ${d.code}] ${d.message}`);

      if (compileResult.success) {
        markdown += 'Result: success\n\n';
      } else {
        markdown += 'Result: failed (invalid bicep)\n\n';
      }
      
      if (diagnostics.length > 0) {
        markdown += 'Diagnostics:\n';
        markdown += `\`\`\`\n${diagnostics.join('\n')}\n\`\`\`\n`;
      }
    } catch (error) {
      markdown += 'Result: failed (unexpected error)\n\n';
      markdown += `\`\`\`\n${error}\n\`\`\`\n`;
    }
  }

  const logPath = path.join(__dirname, 'results.md');
  writeFileSync(logPath, markdown, 'utf8');
  bicep.dispose();
}

series([main]);