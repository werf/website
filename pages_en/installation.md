---
title: Installation
permalink: installation.html
layout: default
sidebar: none
description: How to install werf?
versions:
  - 1.2
  - 1.1
channels:
  - alpha
  - beta
  - ea
  - stable
  - rock-solid
arch:
  - amd64
  - arm64
---
{%- asset installation.css %}
{%- asset installation.js %}
{%- asset releases.css %}

{%- assign releases = site.data.common.releases.releases %}
{%- assign groups = site.data.common.releases_history.history | map: "group" | uniq | reverse %}
{%- assign channels_sorted = site.data.common.channels_info.channels | sort: "stability" %}
{%- assign channels_sorted_reverse = site.data.common.channels_info.channels | sort: "stability" | reverse  %}

<div class="page__container page_installation">

  <div class="installation-selector-row">
    <div class="installation-selector">
      <div id="installation__version" class="installation-selector__title">Version
        <span title="werf follows a versioning strategy called Semantic Versioning. <a href='/development/backward_compatibility.html' class='installation__release-channels--link'>Learn more</a>"></span>
      </div>
      <div class="tabs tabs_simple_condensed">
        <a href="javascript:void(0)" class="tabs__btn"
          data-install-tab-group="version" data-install-tab="1.2">1.2</a>
        <a href="javascript:void(0)" class="tabs__btn"
          data-install-tab-group="version" data-install-tab="1.1">1.1</a>
      </div>
    </div><!-- /selector -->
    <div class="installation-selector">
      <div id="installation__release-channels" class="installation-selector__title">Stability channel
        <span title="All changes in werf go through all release channels. <a href='/development/release_channels.html' class='installation__release-channels--link'>Learn more</a>"></span>
      </div>
      <div class="tabs tabs_simple_condensed">
        <a href="javascript:void(0)" class="tabs__btn"
          data-install-tab-group="channel" data-install-tab="rock-solid">Rock-Solid</a>
        <a href="javascript:void(0)" class="tabs__btn"
          data-install-tab-group="channel" data-install-tab="stable">Stable</a>
        <a href="javascript:void(0)" class="tabs__btn"
          data-install-tab-group="channel" data-install-tab="ea">Early-Access</a>
        <a href="javascript:void(0)" class="tabs__btn"
          data-install-tab-group="channel" data-install-tab="beta">Beta</a>
        <a href="javascript:void(0)" class="tabs__btn"
          data-install-tab-group="channel" data-install-tab="alpha">Alpha</a>
      </div>
    </div><!-- /selector -->
  </div><!-- /selector-row -->
  <div class="installation-selector-row">
    <div class="installation-selector">
      <div class="installation-selector__title">OS</div>
      <div class="tabs tabs_simple_condensed">
        <a href="javascript:void(0)" class="tabs__btn"
          data-install-tab-group="os" data-install-tab="linux">Linux</a>
        <a href="javascript:void(0)" class="tabs__btn"
          data-install-tab-group="os" data-install-tab="macos">Mac OS</a>
        <a href="javascript:void(0)" class="tabs__btn"
          data-install-tab-group="os" data-install-tab="windows">Windows</a>
      </div>
    </div><!-- /selector -->
    <div class="installation-selector">
      <div class="installation-selector__title">Arch</div>
      <div class="tabs tabs_simple_condensed">
        <a href="javascript:void(0)" class="tabs__btn"
          data-install-tab-group="arch" data-install-tab="amd64">Amd64</a>
        <a href="javascript:void(0)" class="tabs__btn"
          data-install-tab-group="arch" data-install-tab="arm64">Arm64</a>
      </div>
    </div><!-- /selector -->
    <div class="installation-selector">
      <div class="installation-selector__title">Installation method</div>
      <div class="tabs tabs_simple_condensed">
        <a href="javascript:void(0)" class="tabs__btn"
          data-install-tab-group="method" data-install-tab="installer">installer</a>
        <a href="javascript:void(0)" class="tabs__btn"
          data-install-tab-group="method" data-install-tab="manually">manually</a>
      </div>
    </div><!-- /selector -->
  </div><!-- /selector-row -->

  <div class="installation-instruction">
      <div class="docs">
<h2 id="install-werf">Installation</h2>
<div class="installation-instruction__tab-content" data-install-content-group="method" data-install-content="manually">
<div class="installation-instruction__tab-content" data-install-content-group="os" data-install-content="linux">
  {% for version in page.versions %}
    <div class="installation-instruction__tab-content" data-install-content-group="version" data-install-content="{{ version }}">
      {% for channel in page.channels %}
        <div class="installation-instruction__tab-content" data-install-content-group="channel" data-install-content="{{ channel }}">
          {% for arch in page.arch %}
            <div class="installation-instruction__tab-content" data-install-content-group="arch" data-install-content="{{ arch }}">
<div markdown="1">
{% include installation/trdl_linux.md version=version channel=channel arch=arch %}
{%- if version != 1.1 %}
{% include installation/setup_buildah.md version=version %}
{%- endif %}
</div>
            </div>
          {% endfor %}
        </div>
      {% endfor %}
    </div>
  {% endfor %}
</div><!-- /os -->
<div class="installation-instruction__tab-content" data-install-content-group="os" data-install-content="macos">
  {% for version in page.versions %}
    <div class="installation-instruction__tab-content" data-install-content-group="version" data-install-content="{{ version }}">
      {% for channel in page.channels %}
        <div class="installation-instruction__tab-content" data-install-content-group="channel" data-install-content="{{ channel }}">
          {% for arch in page.arch %}
            <div class="installation-instruction__tab-content" data-install-content-group="arch" data-install-content="{{ arch }}">
<div markdown="1">{% include installation/trdl_macos.md version=version channel=channel arch=arch %}</div>
            </div>
          {% endfor %}
        </div>
      {% endfor %}
    </div>
  {% endfor %}
</div><!-- /os -->
<div class="installation-instruction__tab-content" data-install-content-group="os" data-install-content="windows">
  {% for version in page.versions %}
    <div class="installation-instruction__tab-content" data-install-content-group="version" data-install-content="{{ version }}">
      {% for channel in page.channels %}
        <div class="installation-instruction__tab-content" data-install-content-group="channel" data-install-content="{{ channel }}">
          {% for arch in page.arch %}
            <div class="installation-instruction__tab-content" data-install-content-group="arch" data-install-content="{{ arch }}">
<div markdown="1">{% include installation/trdl_windows.md version=version channel=channel arch=arch %}</div>
            </div>
          {% endfor %}
        </div>
      {% endfor %}
    </div>
  {% endfor %}
</div><!-- /os -->

      </div><!-- /method -->
      <div class="installation-instruction__tab-content" data-install-content-group="method" data-install-content="installer">
        <div class="installation-instruction__tab-content" data-install-content-group="os" data-install-content="linux">
  {% for version in page.versions %}
    <div class="installation-instruction__tab-content" data-install-content-group="version" data-install-content="{{ version }}">
      {% for channel in page.channels %}
        <div class="installation-instruction__tab-content" data-install-content-group="channel" data-install-content="{{ channel }}">
          {% for arch in page.arch %}
            <div class="installation-instruction__tab-content" data-install-content-group="arch" data-install-content="{{ arch }}">
<div markdown="1">
{% include installation/installer_linux_macos.md version=version channel=channel %}
{%- if version != 1.1 %}
{% include installation/setup_buildah.md version=version %}
{%- endif %}
</div>
            </div>
          {% endfor %}
        </div>
      {% endfor %}
    </div>
  {% endfor %}

        </div>
        <div class="installation-instruction__tab-content" data-install-content-group="os" data-install-content="macos">
  {% for version in page.versions %}
    <div class="installation-instruction__tab-content" data-install-content-group="version" data-install-content="{{ version }}">
      {% for channel in page.channels %}
        <div class="installation-instruction__tab-content" data-install-content-group="channel" data-install-content="{{ channel }}">
          {% for arch in page.arch %}
            <div class="installation-instruction__tab-content" data-install-content-group="arch" data-install-content="{{ arch }}">
<div markdown="1">
{% include installation/installer_linux_macos.md version=version channel=channel %}
</div>
            </div>
          {% endfor %}
        </div>
      {% endfor %}
    </div>
  {% endfor %}
        </div>
      </div><!-- /method -->
    </div>
  </div>

  
</div>
