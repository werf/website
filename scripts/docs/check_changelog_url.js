const assert = require("assert");
const fs = require("fs");
const path = require("path");

const src = fs.readFileSync(
  path.join(__dirname, "../../assets/js/changelog.js"),
  "utf8",
);
const fnSource = src.match(
  /function getChangelogURL\(href, version, channel\) \{[\s\S]*?\n\}/,
);
assert.ok(fnSource, "getChangelogURL is not found in changelog.js");

const getChangelogURL = new Function(
  `${fnSource[0]}; return getChangelogURL;`,
)();

const changelogURL = "https://werf.io/about/changelog.html";

const cases = [
  [
    "release and release channel are added",
    [changelogURL, "2", "stable"],
    `${changelogURL}?version=2&channel=stable`,
  ],
  [
    "release and release channel are replaced",
    [`${changelogURL}?version=2&channel=stable`, "1.2", "rock-solid"],
    `${changelogURL}?version=1.2&channel=rock-solid`,
  ],
  [
    "other query parameters and the fragment are kept",
    [`${changelogURL}?utm_source=rss&channel=ea#releases`, "2", "alpha"],
    `${changelogURL}?utm_source=rss&channel=alpha&version=2#releases`,
  ],
  [
    "release channel is removed if no channel is open",
    [`${changelogURL}?version=2&channel=stable`, "1.2", undefined],
    `${changelogURL}?version=1.2`,
  ],
  [
    "values are encoded",
    [changelogURL, "a&b", "c d"],
    `${changelogURL}?version=a%26b&channel=c+d`,
  ],
];

for (const [name, args, want] of cases) {
  assert.strictEqual(getChangelogURL(...args), want, name);
}

console.log(`OK: getChangelogURL, ${cases.length} cases`);
