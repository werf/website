---
title: "Kubernetes + werf tutorial: learn K8s and deploy your app"
layout: guides-landing
permalink: guides/index.html
---

<div class="landing">
    <div class="landing__header">
        <div class="landing__container">
            <a href="/" class="landing__header-title" data-proofer-ignore>
                {% asset werf-logo.svg alt="werf" width="167" height="70" %}
            </a>
            <a href="{{ site.site_urls['en'] }}" class="landing__button" data-proofer-ignore>
                {% asset arrow.svg %}
                <span>werf.io</span>
            </a>
        </div>
    </div>
    <div class="landing__content">
        <div class="landing__container">
            <div class="landing__section landing__section_first" data-sm-trigger="intro">
                <h1 class="landing__h1">
                    Delivering to&nbsp;Kubernetes<br>
                    quickly and&nbsp;effortlessly<br>
                    using&nbsp;<b>werf</b>
                </h1>
                <div class="landing__text">
                    In this free tutorial, you will learn how to&nbsp;make your applications ready for Kubernetes
                    and continuously deploy them by implementing CI/CD using Open Source tools.
                </div>
            </div>
            <div class="landing__section" data-sm-trigger="plan">
                <h1 class="landing__h2">
                    What am I getting into?
                </h1>
                <div class="landing__text">
                    <p>This tutorial combines the <b>theory and practice</b> of development (Dev) and operation (Ops).</p>
                    <p>Its contents are aimed at developers seeking to acquire basic DevOps skills in organizing the continuous delivery of applications to Kubernetes. The DevOps engineers who want to solve their tasks more efficiently will also benefit from it.</p>
                    <p>We will gradually cover <b>Kubernetes basics</b>, <b>preparing your applications to run in Kubernetes</b> and all the tasks related to&nbsp;developing services and implementing the CI/CD process for them (building, deploying, working with dependencies and&nbsp;assets, working with&nbsp;databases and&nbsp;files, storing sensitive and non-sensitive application configurations) as&nbsp;well&nbsp;as the <b>best practices for&nbsp;deploying</b> in&nbsp;K8s.</p>
                    <p>The guides take into&nbsp;account the&nbsp;specifics of&nbsp;programming languages/frameworks and include examples of the&nbsp;application source&nbsp;code and related infrastructure (IaC).</p>
                    <p><i>This tutorial is based on&nbsp;a&nbsp;plain "vanilla" Kubernetes cluster. However, you can easily adapt it for&nbsp;custom platforms.</i></p>
                </div>
            </div>
            <div class="landing__section" data-sm-trigger="todo">
                <h1 class="landing__h2">
                    What's in it for me?
                </h1>
                <div class="landing__text">
                    <p>Mastering the skills using a self-study guide is a difficult challenge. You have to find the time, work hard, and focus on the process.</p>
                    <p>If you go all the way to the end you will be able to:</p>
                </div>
                <ul class="landing__list">
                    <li>
                        Gain expertise and better market your skills to employers
                        <span>(in the area that combines the development and operation and commonly referred to as DevOps)</span>
                    </li>
                    <li>
                        Introduce Kubernetes to your company
                    </li>
                    <li>
                        Learn how to solve problems in your company more efficiently
                    </li>
                </ul>
                <div class="landing__text">
                    <p>This tutorial will help you focus on your needs, and the growing <a href="https://t.me/werf_io">werf user community</a> will make it easier to overcome any obstacles along the way.</p>
                </div>
            </div>
            <div class="landing__section" data-sm-trigger="select">
                <h1 class="landing__h2">
                    Dive in! <span>(werf 1.2)</span>
                </h1>

                <div class="landing__text">
                  <p>Choose the framework that fits you best:</p>
                </div>

                {% include common/guides-landing-tiles.html %}
            </div>
            <div class="landing__section">
            </div>
        </div>
    </div>
    <div class="landing__side">
        <div class="landing__captain" data-sm-captain="">
            {% asset landing/parrot.png class="landing__captain-img landing__captain-img_intro" %}
            {% asset landing/parrot2.png class="landing__captain-img landing__captain-img_plan" %}
            {% asset landing/parrot3.png class="landing__captain-img landing__captain-img_todo" %}
            {% asset landing/parrot4.png class="landing__captain-img landing__captain-img_select" %}
        </div>
    </div>
</div>
