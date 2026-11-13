const assert = require("assert");
const fs = require("fs");
const path = require("path");

const src = fs.readFileSync(
  path.join(__dirname, "../../assets/js/report-typo.js"),
  "utf8",
);
const fnSource = src.match(
  /function getTypoReportURL\(issueURL, location, selection\) \{[\s\S]*?\n\}/,
);
assert.ok(fnSource, "getTypoReportURL is not found in report-typo.js");

const getTypoReportURL = new Function(
  `${fnSource[0]}; return getTypoReportURL;`,
)();

const issueURL = "https://github.com/werf/website/issues/new";
const location = {
  href: "https://werf.io/guides/nodejs/100_basic/10_build.html?a=b#build",
  pathname: "/guides/nodejs/100_basic/10_build.html",
};

const report = (selection) => {
  const url = new URL(getTypoReportURL(issueURL, location, selection));
  assert.strictEqual(url.origin + url.pathname, issueURL);
  return {
    title: url.searchParams.get("title"),
    body: url.searchParams.get("body"),
  };
};

const reportBody = (fence, text) =>
  [
    "Page: https://werf.io/guides/nodejs/100_basic/10_build.html?a=b",
    "",
    "Text with the error:",
    "",
    fence,
    text,
    fence,
    "",
    "<!-- Describe the error and suggest a fix if you can. -->",
    "",
  ].join("\n");

const cases = [
  ["plain text", "Teh image", reportBody("```", "Teh image")],
  [
    "surrounding whitespace is trimmed",
    "\n  Teh image \t\n",
    reportBody("```", "Teh image"),
  ],
  [
    "line breaks are normalized",
    "first line\r\nsecond line\rthird line",
    reportBody("```", "first line\nsecond line\nthird line"),
  ],
  [
    "Markdown, mentions and issue references are kept as is",
    "**werf** @werf #123 > quote",
    reportBody("```", "**werf** @werf #123 > quote"),
  ],
  [
    "fence is longer than any backtick run in the text",
    "run `werf build` or ```werf converge```",
    reportBody("````", "run `werf build` or ```werf converge```"),
  ],
  [
    "long text is truncated",
    "a".repeat(1500),
    reportBody("```", "a".repeat(1000) + "…"),
  ],
];

for (const [name, selection, want] of cases) {
  const got = report(selection);
  assert.strictEqual(
    got.title,
    "Typo on /guides/nodejs/100_basic/10_build.html",
    name,
  );
  assert.strictEqual(got.body, want, name);
}

console.log(`OK: getTypoReportURL, ${cases.length} cases`);
