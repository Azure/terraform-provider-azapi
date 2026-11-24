from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]


def _normalized(name: str) -> str:
    return " ".join((ROOT / "definitions" / name).read_text(encoding="utf-8").split())


def test_azapi_follow_up_acknowledges_supplied_evidence():
    charter = _normalized("tf_requirements.md")

    required_guidance = (
        "inventory the concrete evidence already supplied",
        "exact resolved version from `.terraform.lock.hcl`",
        "plan showing a persistent diff is evidence that the drift exists",
        "summary of the relevant evidence already received",
        "Never re-request a resource type, API version, configuration block",
        "Do not say that the creator's reply \"doesn't answer\"",
    )

    for guidance in required_guidance:
        assert guidance in charter


def test_requirements_charter_owns_precise_collaborative_follow_ups():
    charter = _normalized("tf_requirements.md")

    required_guidance = (
        "Begin by summarizing the concrete evidence already supplied",
        "ask specifically for the resolved/locked version",
        "Distinguish evidence that the symptom exists",
        "do not say that the creator's reply \"doesn't answer\"",
    )

    for guidance in required_guidance:
        assert guidance in charter
