# BACKLOG

- use bubbletea to make resumake a lightweight tui ux
- add progress indicators/spinners, especially during the api call
- provide more verbose logging/debugging options (`--verbose`)
- implement a more robust stream-of-consciousness capture (e.g., opening `$EDITOR`). should also support quite extensive input; we're working with a model that supports over a million tokens of context
- should be allowed to include an arbitrary number of text/markdown files as context
- add *sanity check* step after initial generation, to make sure nothing in the generated resume is made up, no fake numbers, etc; all specifics must be sourced in the context provided; then generate a second draft with the sanity check's added context
- add second draft step, so we generate a first draft then run it through a critic, then generate a second draft with the first draft and the critique and all original context
- support flow for fine-tuning resume given a particular job posting / description
- generate cover letters when generating resume fine tuned on a particular job posting
- integrate with github to pull contributions / activity as context
