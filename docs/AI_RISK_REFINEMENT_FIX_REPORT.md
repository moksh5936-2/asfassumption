# AI Risk Refinement Fix Report

## Problem

`parseRiskRefinements()` at `asf-tui/ai.go:160` returned `nil` unconditionally. The AI prompt asked the LLM for risk refinements, and the response parsing code in `parseResponse` dispatched to this stub, but no actual parsing was ever implemented. AI-powered risk refinement was advertised in the prompt and README but was non-functional.

## Decision: Option B — Remove the claim and the prompt section

Since implementing real risk refinement parsing requires:
1. A defined structured output format from the LLM
2. A reliable parsing function that extracts `AIRiskRefinement` structs (AssumptionID, OriginalRisk, SuggestedRisk, Reasoning)

…and the current prompt does not constrain the LLM to produce parseable risk refinements, we chose Option B: remove the prompt section and the README claim rather than ship fake functionality.

## Changes Made

### `asf-tui/ai.go`

1. **Prompt section removed** (line 92): The `2. RISK_REFINEMENTS: Any assumptions where the risk level seems incorrect.` line was deleted from `buildPrompt()`. The remaining items were renumbered (3→2, 4→3).

2. **`parseRiskRefinements` documented** (line 160): Added a doc comment explaining the function is reserved for future use and what would be needed for a real implementation. The function body remains `return nil`.

3. **`parseResponse` left unchanged**: The `"risks"` case at line 116 still handles LLM output containing "risk_refinement" or "risk refinement". Since the prompt no longer asks for this section, the LLM will never emit it, so the dead route is harmless. The `append` of nil is a no-op.

### `README.md`

1. **Line 294**: Changed "additional assumptions, risk refinements, and recommendations" → "additional assumptions and recommendations".

2. **Limitations section**: Added a new limitation entry documenting that AI risk refinement is not yet implemented, with a cross-reference to this report.

## Future Recommendation

To re-enable risk refinement in the future:

1. Define a structured output format in the prompt, e.g.:
   ```
   RISK_REFINEMENTS:
   - AssumptionID: "1"
     SuggestedRisk: "High"
     Reasoning: "The assumption lacks encryption, elevating risk."
   ```
2. Implement `parseRiskRefinements` to parse this format using regex or line-by-line parsing.
3. Implement `mergeAIResults` to consume `enhanced.RefinedRisks` and update original assumption risk levels.
4. Remove the limitation note from README.
5. Re-add the RISK_REFINEMENTS section to the prompt.
