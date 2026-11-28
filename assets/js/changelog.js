// Changelog page: keep the selected release and release channel tabs in the URL query
// (e.g., /about/changelog.html?version=2&channel=stable), so that a shared or bookmarked link opens the same tabs.

function getChangelogURL(href, version, channel) {
  var url = new URL(href);

  url.searchParams.set('version', version);
  if (channel) {
    url.searchParams.set('channel', channel);
  } else {
    url.searchParams.delete('channel');
  }

  return url.toString();
}

$(document).ready(function () {
  var $versionButtons = $('[data-changelog-version]:not([data-changelog-channel])');
  var $channelButtons = $('[data-changelog-channel]');

  // Compare the attribute values instead of building a selector from the URL query, which can contain anything.
  function filterButtons($buttons, attribute, value) {
    return $buttons.filter(function () {
      return this.getAttribute(attribute) === value;
    });
  }

  function getChannelButtons(version) {
    return filterButtons($channelButtons, 'data-changelog-version', version);
  }

  // Open the tabs requested in the URL query. Unknown values are ignored, and the URL is left as is.
  var params = new URLSearchParams(window.location.search);

  var versionButton = filterButtons($versionButtons, 'data-changelog-version', params.get('version')).get(0);
  if (versionButton) {
    versionButton.click();
  }

  var activeVersion = $versionButtons.filter('.active').attr('data-changelog-version');
  var channelButton = filterButtons(getChannelButtons(activeVersion), 'data-changelog-channel', params.get('channel')).get(0);
  if (channelButton) {
    channelButton.click();
  }

  // Selecting a release keeps the release channel that is open within it.
  $versionButtons.add($channelButtons).on('click', function () {
    var version = this.getAttribute('data-changelog-version');
    var channel = this.getAttribute('data-changelog-channel') || getChannelButtons(version).filter('.active').attr('data-changelog-channel');

    window.history.replaceState(null, null, getChangelogURL(window.location.href, version, channel));
  });
});
