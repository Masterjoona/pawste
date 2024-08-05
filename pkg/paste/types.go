package paste

import (
	"mime/multipart"
)

type File struct {
	ID          int
	Name        string
	Size        int
	ContentType string
	Blob        []byte
}

type Paste struct {
	ID          int
	PasteName   string
	Expire      int64
	Privacy     string
	NeedsAuth   int
	ReadCount   int
	ReadLast    int64
	BurnAfter   int
	Content     string
	UrlRedirect int
	Syntax      string
	Password    string
	Files       []File
	CreatedAt   int64
	UpdatedAt   int64
}

type PasteLists struct {
	Pastes    []Paste
	Redirects []Paste
}

type PasteUpdate struct {
	Content        string                  `form:"content,omitempty"`
	Password       string                  `form:"password,omitempty"`
	RemovedFiles   []string                `form:"removed_files,omitempty"`
	FilesMultiPart []*multipart.FileHeader `form:"file,omitempty"`
	Files          []File
}

type LookUpMap struct {
	Items map[string]struct{}
}

func createLookupMap(slice []string) *LookUpMap {
	lookup := make(map[string]struct{}, len(slice))
	for _, item := range slice {
		lookup[item] = struct{}{}
	}
	return &LookUpMap{Items: lookup}
}

func (l *LookUpMap) Contains(str string) bool {
	_, found := l.Items[str]
	return found
}

var privacyOptions = []string{"public", "unlisted", "readonly", "private", "secret"}

var syntaxOptions = []string{
	"none",
	"abap",
	"actionscript-3",
	"ada",
	"angular-html",
	"angular-ts",
	"apache",
	"apex",
	"apl",
	"applescript",
	"ara",
	"asciidoc",
	"asm",
	"astro",
	"awk",
	"ballerina",
	"bat",
	"beancount",
	"berry",
	"bibtex",
	"bicep",
	"blade",
	"c",
	"cadence",
	"clarity",
	"clojure",
	"cmake",
	"cobol",
	"codeowners",
	"codeql",
	"coffee",
	"common-lisp",
	"cpp",
	"crystal",
	"csharp",
	"css",
	"csv",
	"cue",
	"cypher",
	"d",
	"dart",
	"dax",
	"desktop",
	"diff",
	"docker",
	"dotenv",
	"dream-maker",
	"edge",
	"elixir",
	"elm",
	"emacs-lisp",
	"erb",
	"erlang",
	"fennel",
	"fish",
	"fluent",
	"fortran-fixed-form",
	"fortran-free-form",
	"fsharp",
	"gdresource",
	"gdscript",
	"gdshader",
	"genie",
	"gherkin",
	"git-commit",
	"git-rebase",
	"gleam",
	"glimmer-js",
	"glimmer-ts",
	"glsl",
	"gnuplot",
	"go",
	"graphql",
	"groovy",
	"hack",
	"haml",
	"handlebars",
	"haskell",
	"haxe",
	"hcl",
	"hjson",
	"hlsl",
	"html",
	"html-derivative",
	"http",
	"hxml",
	"hy",
	"imba",
	"ini",
	"java",
	"javascript",
	"jinja",
	"jison",
	"json",
	"json5",
	"jsonc",
	"jsonl",
	"jsonnet",
	"jssm",
	"jsx",
	"julia",
	"kotlin",
	"kusto",
	"latex",
	"lean",
	"less",
	"liquid",
	"log",
	"logo",
	"lua",
	"luau",
	"make",
	"markdown",
	"marko",
	"matlab",
	"mdc",
	"mdx",
	"mermaid",
	"mojo",
	"move",
	"narrat",
	"nextflow",
	"nginx",
	"nim",
	"nix",
	"nushell",
	"objective-c",
	"objective-cpp",
	"ocaml",
	"pascal",
	"perl",
	"php",
	"plsql",
	"po",
	"postcss",
	"powerquery",
	"powershell",
	"prisma",
	"prolog",
	"proto",
	"pug",
	"puppet",
	"purescript",
	"python",
	"qml",
	"qmldir",
	"qss",
	"r",
	"racket",
	"raku",
	"razor",
	"reg",
	"regexp",
	"rel",
	"riscv",
	"rst",
	"ruby",
	"rust",
	"sas",
	"sass",
	"scala",
	"scheme",
	"scss",
	"shaderlab",
	"shellscript",
	"shellsession",
	"smalltalk",
	"solidity",
	"soy",
	"sparql",
	"splunk",
	"sql",
	"ssh-config",
	"stata",
	"stylus",
	"svelte",
	"swift",
	"system-verilog",
	"systemd",
	"tasl",
	"tcl",
	"templ",
	"terraform",
	"tex",
	"toml",
	"ts-tags",
	"tsv",
	"tsx",
	"turtle",
	"twig",
	"typescript",
	"typespec",
	"typst",
	"v",
	"vala",
	"vb",
	"verilog",
	"vhdl",
	"viml",
	"vue",
	"vue-html",
	"vyper",
	"wasm",
	"wenyan",
	"wgsl",
	"wikitext",
	"wolfram",
	"xml",
	"xsl",
	"yaml",
	"zenscript",
	"zig",
}

var PrivacyMap = createLookupMap(privacyOptions)
var SyntaxMap = createLookupMap(syntaxOptions)
