$(document).ready(function () {
  var default_os;
  if (bowser.windows)
    default_os = 'windows'
  else if (bowser.mac)
    default_os = 'macos'
  else
    default_os = 'linux'

  var defaults = {
    version: '2',
    channel: 'stable',
    os: default_os,
    arch: 'amd64'
  }

  const windows_default_method = "manually";
  const unix_default_method = "installer";

  function doInstallSelect(group, param) {
    $(`[data-install-tab-group="${group}"]`).removeClass('active');
    $(`[data-install-tab="${param}"]`).addClass('active');

    $(`[data-install-content-group="${group}"]`).removeClass('active');
    $(`[data-install-content="${param}"]`).addClass('active');

    $(`[data-install-info="${group}"]`).text($(`[data-install-tab="${param}"]`).text());
  }

  function installSelect(group, param) {
    // Update URL params
    let url = new URL(window.location.href);
    let params = new URLSearchParams(url.search);
    params.set(group, param);
    url.search = params.toString();
    window.history.replaceState(null, null, url.toString());

    // Update buttons status
    if (group == "version") {
      if (param == "2") {
        $(`[data-install-tab="rock-solid"]`).hide();
      } else {
        $(`[data-install-tab="rock-solid"]`).show();
      }
    }

    if (group == "os") {
      if (param == "windows") {
        $(`[data-install-tab="arm64"]`).hide();
        doInstallSelect("arch", "amd64")

        $(`[data-install-tab="installer"]`).hide();
        doInstallSelect("method", windows_default_method);
      } else {
        $(`[data-install-tab="arm64"]`).show();
        $(`[data-install-tab="amd64"]`).show();
        $(`[data-install-tab="installer"]`).show();
      }
    }

    doInstallSelect(group, param)
  }

  let url = new URL(window.location.href);
  let params = new URLSearchParams(url.search);
  Object.keys(defaults).forEach(function (key) {
    if (!params.get(key)) {
      params.set(key, defaults[key]);
    }
  });

  for (let [key, value] of params) {
    installSelect(key, value)
  }
  if (!params.get("method")) {
    if (params.get("os") === "windows") {
      installSelect("method", windows_default_method)
    } else {
      installSelect("method", unix_default_method)
    }
  }

  $('[data-install-tab]').on('click', function () {
    installSelect($(this).data('install-tab-group'), $(this).data('install-tab'));
  })

})
