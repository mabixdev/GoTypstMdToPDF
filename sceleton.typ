#import "@preview/mitex:0.2.4": mitex
#import "@preview/cmarker:0.1.1"

// Document setup
#set page(margin: 2cm)
#set text(font: "Arial", size: 11pt)
#set par(justify: true, leading: 0.65em)

// Math and code styling
#show math.equation: set block(spacing: 1em)
#set math.equation(numbering: "(1)")

// Header styling
#show heading.where(level: 1): it => [
  #set text(size: 16pt, weight: "bold")
  #block(spacing: 1.5em)[#it]
]

#show heading.where(level: 2): it => [
  #set text(size: 14pt, weight: "bold")
  #block(spacing: 1.2em)[#it]
]

#show heading.where(level: 3): it => [
  #set text(size: 12pt, weight: "bold")
  #block(spacing: 1em)[#it]
]

#cmarker.render(`
{{Placeholder Markdown}}
`, math: mitex)