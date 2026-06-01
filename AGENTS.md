## Core principles

- Stop telling me "You're right", it just shows how incompetent you are. Do it right on your first try, fact-check and review after changes. If you are not sure, ask for help.
- When you see changes made outside your knowledge, use the current version as your new starting point. Do not blindly overwrite those changes or you suck. Even if you have to update the code, always respect the pattern in the surrounding context!

## Style and mechanics

This applies to all texts, including but not limited to UI, documentation, code comments.

- Use sentence case. Preserve original casing for brand names.
- End with a period for a full sentence.
- Never use em dashes (`—`) or en dashes (`–`) in prose. Rewrite the sentence with a comma, period, colon, or parentheses instead. Exception: em/en dashes are allowed as visual separators in UI design (e.g., between a title and description, in a terminal prompt label) where they function as a graphic element rather than punctuation.
- Do not overuse semicolons. Two short sentences are almost always clearer than one sentence joined by a semicolon. Reserve the semicolon for the rare case where the two clauses are so tightly coupled that splitting them loses meaning, never as a default em-dash replacement or a way to chain related thoughts.
- Do not add comments that repeat what the code is doing, always prefer more descriptive names. Do add comments for intentions that aren't obvious via reading the code alone. This rule takes precedence over matching existing patterns.

## Coding guidelines

- Use `github.com/cockroachdb/errors` for error handling.
- Use `github.com/stretchr/testify` for assertions in tests. Be mindful about the choice of `require` and `assert`, the former should be used when the test cannot proceed meaningfully after a failed assertion.
- Always use backtick (raw string) literals for multi-line strings. Never use `"...\n" +` concatenation.
- Prefer grouping related cases as `t.Run` subtests under a single top-level `TestXxx` function, instead of writing one top-level `TestXxx_Case` per case. Reserve separate top-level functions for unrelated subjects and complex subtest setup.
- Always use backtick (raw string) literals for multi-line strings. Never use `"...\n" +` concatenation.

## Source code control

- When pushing changes to a pull request from a fork, use SSH address and do not add remote.
- Never commit on the `main` branch directly unless being explicitly asked to do so. A single ask only grants a single commit action on the `main` branch.
- Never amend commits unless being explicitly asked to do so.
- When creating a git worktree, the worktree directory name must match its branch name. Do not use random or generated suffixes.
