// Report a typo or an error found on the page: select the text and press Ctrl+Enter (Cmd+Enter on macOS)
// to open a pre-filled GitHub issue. Enabled only on pages that render the footer hint with the
// `data-report-typo` attribute, which holds the URL of the "new issue" page.

function getTypoReportURL(issueURL, location, selection) {
  var maxSelectionLength = 1000;
  var text = selection.replace(/\r\n?/g, '\n').trim();

  // Truncate by code points, so that a character outside the BMP (e.g., an emoji) is never cut in half:
  // a lone surrogate makes encodeURIComponent() throw.
  var chars = Array.from(text);
  if (chars.length > maxSelectionLength) {
    text = chars.slice(0, maxSelectionLength).join('').trimEnd() + '…';
  }

  // Put the text into a fenced code block to keep it as is: no Markdown rendering, mentions or issue references.
  var fenceLength = (text.match(/`+/g) || []).reduce(function (length, backticks) {
    return Math.max(length, backticks.length + 1);
  }, 3);
  var fence = new Array(fenceLength + 1).join('`');

  var title = 'Typo on ' + location.pathname;
  var body = [
    'Page: ' + location.href.split('#')[0],
    '',
    'Text with the error:',
    '',
    fence,
    text,
    fence,
    '',
    '<!-- Describe the error and suggest a fix if you can. -->',
    ''
  ].join('\n');

  return issueURL + '?title=' + encodeURIComponent(title) + '&body=' + encodeURIComponent(body);
}

$(document).ready(function () {
  var issueURL = $('[data-report-typo]').data('report-typo');

  if (!issueURL) {
    return;
  }

  $(document).on('keydown', function (e) {
    if (e.key !== 'Enter' || !(e.ctrlKey || e.metaKey) || e.altKey || e.shiftKey || (e.originalEvent && e.originalEvent.isComposing)) {
      return;
    }

    // Do not interfere with Ctrl+Enter in form fields and editable content.
    if ($(e.target).closest('input, textarea, select, [contenteditable]').length) {
      return;
    }

    var selection = window.getSelection().toString();
    if (!selection.trim()) {
      return;
    }

    e.preventDefault();
    window.open(getTypoReportURL(issueURL, window.location, selection), '_blank', 'noopener');
  });
});
