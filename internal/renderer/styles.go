package renderer

// CustomStyle returns a custom Glamour style JSON without heading prefixes
// This provides a cleaner look similar to modern markdown viewers
const CustomStyle = `{
  "document": {
    "block_prefix": "\n",
    "block_suffix": "\n",
    "color": "252",
    "margin": 2
  },
  "block_quote": {
    "indent": 1,
    "indent_token": "│ "
  },
  "paragraph": {},
  "list": {
    "level_indent": 2
  },
  "heading": {
    "block_suffix": "\n",
    "color": "39",
    "bold": true
  },
  "h1": {
    "prefix": " ",
    "suffix": " ",
    "color": "228",
    "background_color": "63",
    "bold": true,
    "block_suffix": "\n"
  },
  "h2": {
    "prefix": " ",
    "color": "39",
    "bold": true,
    "block_suffix": "\n"
  },
  "h3": {
    "prefix": " ",
    "color": "41",
    "bold": true
  },
  "h4": {
    "prefix": "  ",
    "color": "42",
    "bold": true
  },
  "h5": {
    "prefix": "  ",
    "color": "43",
    "bold": false
  },
  "h6": {
    "prefix": "  ",
    "color": "44",
    "bold": false,
    "italic": true
  },
  "strikethrough": {
    "crossed_out": true
  },
  "emph": {
    "italic": true
  },
  "strong": {
    "bold": true
  },
  "hr": {
    "color": "240",
    "format": "\n────────────────────────────────────────────────────────────────────────────────\n"
  },
  "item": {
    "block_prefix": "• "
  },
  "enumeration": {
    "block_prefix": ". "
  },
  "task": {
    "ticked": "[✓] ",
    "unticked": "[ ] "
  },
  "link": {
    "color": "30",
    "underline": true
  },
  "link_text": {
    "color": "35",
    "bold": true
  },
  "image": {
    "color": "212",
    "underline": true
  },
  "image_text": {
    "color": "243",
    "format": "Image: {{.text}} →"
  },
  "code": {
    "prefix": " ",
    "suffix": " ",
    "color": "203",
    "background_color": "236"
  },
  "code_block": {
    "color": "244",
    "margin": 2,
    "chroma": {
      "text": {
        "color": "#C4C4C4"
      },
      "error": {
        "color": "#F1F1F1",
        "background_color": "#F05B5B"
      },
      "comment": {
        "color": "#676767"
      },
      "comment_preproc": {
        "color": "#FF875F"
      },
      "keyword": {
        "color": "#00AAFF"
      },
      "keyword_reserved": {
        "color": "#FF5FD2"
      },
      "keyword_namespace": {
        "color": "#FF5F87"
      },
      "keyword_type": {
        "color": "#6E88FF"
      },
      "operator": {
        "color": "#EF8080"
      },
      "punctuation": {
        "color": "#E8E8A8"
      },
      "name": {
        "color": "#C4C4C4"
      },
      "name_builtin": {
        "color": "#FF8EC7"
      },
      "name_tag": {
        "color": "#B083EA"
      },
      "name_attribute": {
        "color": "#7A82DA"
      },
      "name_class": {
        "color": "#F1F1F1",
        "underline": true,
        "bold": true
      },
      "name_constant": {
        "color": "#FF5FD2"
      },
      "name_decorator": {
        "color": "#FFFF87"
      },
      "name_function": {
        "color": "#00D787"
      },
      "literal_number": {
        "color": "#6EEFC0"
      },
      "literal_string": {
        "color": "#C69669"
      },
      "literal_string_escape": {
        "color": "#AFFFD7"
      },
      "generic_deleted": {
        "color": "#FD5B5B"
      },
      "generic_emph": {
        "italic": true
      },
      "generic_inserted": {
        "color": "#00D787"
      },
      "generic_strong": {
        "bold": true
      },
      "generic_subheading": {
        "color": "#777777"
      },
      "background": {
        "background_color": "#373737"
      }
    }
  },
  "table": {
    "center_separator": "┼",
    "column_separator": "│",
    "row_separator": "─"
  },
  "definition_list": {},
  "definition_term": {},
  "definition_description": {
    "block_prefix": "\n🠶 "
  },
  "html_block": {},
  "html_span": {}
}`
