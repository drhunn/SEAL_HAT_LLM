# NOTE

## cli.py retirement
The executable package entrypoint no longer relies on `cli.py`.

The canonical executable entrypoint is now:
- `python/hat_llm/build_dataset.py`

The old `cli.py` implementation has been preserved as documentation in:
- `python/hat_llm/cli.md`

## why
This change keeps the runnable entrypoint clear while preserving the previous implementation in Markdown form for reference.

## standing convention
If a future in-session update cannot safely land as a live `.py` module, it may be preserved as a `.md` file in this folder with an explicit note like this one.
