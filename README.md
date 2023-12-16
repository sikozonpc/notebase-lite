# Notebase

(WIP) A centralized highlights and notes ecosystem that supports multiple sources (Kindle, Notion, etc.)

## Features

- Automatically import from Kindle and other sources.
- Notes are automatically categorized and indexed.
- Daily review of notes by email or REST.
- Suggested content based on your notes.
- Advanced search and filtering.
- Export to Markdown, HTML, etc.

## Installation

Make sure to have GO 1.12+ installed. And then run:
```bash
make run
```

## Technical

The project is built just using the standard library and a few essentials (MySQL driver and mux). The project is built in a such a way
to be as minimal and simple to test how powerful the GO standard library goes. However as the project grows it might be subject to change.