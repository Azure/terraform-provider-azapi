<!DOCTYPE html>
<html>
<head />
<body>
    <p style="font-size:12pt;font-family:'Segoe UI'">Dear $DevFirstName $DevLastName,</p>
    <p style="font-size:12pt;font-family:'Segoe UI'">
        Your account has been created. Please follow the link below to visit the $OrganizationName developer portal
        and claim it:
    </p>
    <p style="font-size:12pt;font-family:'Segoe UI'">
        <a
            href="${developer_portal_url}/confirm-v2/identities/basic/invite?$ConfirmQuery">${developer_portal_url}/confirm-v2/identities/basic/invite?$ConfirmQuery</a>
    </p>
    <p style="font-size:12pt;font-family:'Segoe UI'">Best,</p>
    <p style="font-size:12pt;font-family:'Segoe UI'">The $OrganizationName API Team</p>
</body>
</html>