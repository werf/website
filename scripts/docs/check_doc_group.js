const assert = require("assert");
const fs = require("fs");
const path = require("path");

const src = fs.readFileSync(
  path.join(__dirname, "../../assets/js/customscripts.js"),
  "utf8",
);
const fnSource = src.match(/function getDocGroupFromURL\(\) \{[\s\S]*?\n\}/);
assert.ok(fnSource, "getDocGroupFromURL is not found in customscripts.js");

const call = (pathname) =>
  new Function("window", `${fnSource[0]}; return getDocGroupFromURL();`)({
    location: { pathname },
  });

const cases = [
  ["/docs/latest/usage/build/overview.html", "v3"],
  ["/docs/pr-123/index.html", "v3"],
  ["/docs/v3/", "v3"],
  ["/docs/v3.1.4/reference/werf_yaml.html", "v3"],
  ["/docs/v2/", "v2"],
  ["/docs/v2.10.3/reference/werf_yaml.html", "v2"],
  ["/docs/v1.2/", "v1.2"],
  ["/docs/v1.2.4-plus-fix1/reference/werf_yaml.html", "v1.2"],
  ["/docs/v9/", null],
  ["/guides.html", null],
];

for (const [pathname, want] of cases) {
  const got = call(pathname);
  assert.strictEqual(got, want, `${pathname}: got ${got}, want ${want}`);
}

console.log(`OK: getDocGroupFromURL, ${cases.length} cases`);
