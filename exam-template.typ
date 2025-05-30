#import "@preview/ttt-exam:0.1.2": *
#import "@preview/mitex:0.2.4": mitex
#import "@preview/cmarker:0.1.1"

#set text(size: 12pt, font: ("Arial"), weight: 400, lang: "en")

// #show: exam.with(
//   class: "Examination",
//   subject: "Academic Subject",
//   title: "Exam",
//   subtitle: "",
//   date: datetime.today().display("[day].[month].[year]"),
//   authors: "Instructor",
//   logo: none,
//   cover: false,
//   eval-table: false,
//   appendix: none,
// )

// Render the markdown content within the exam structure
#cmarker.render(`
{{Placeholder Markdown}}
`, math: mitex) 