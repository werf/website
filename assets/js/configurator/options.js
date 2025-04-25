$(document).ready(() => {
    let initialData = {};
    let optionsList;
    let pathToData;
    const DATA_FILE = '/configurator-data.json';
    const OPTIONS_LIST_FILE = '/configurator-options-list.json';

    // Load options state from the query string.
    if (window.location.search) {
        let urlParams = window.location.search.slice(1);
        const params = urlParams.split('&');
        const initPathToData = params.join('_').replaceAll('=', '-');

        params.forEach((parameter) => {
            const paramCortege = parameter.split('=');
            initialData[paramCortege[0]] = paramCortege[1];
        })

        getData(DATA_FILE).then(res => {
            if (res[initPathToData]) {
                $.ajax({
                    type: "GET",
                    url: `${res[initPathToData]}`,
                    success: function (response) {
                        $('#configurator-content').html(response);
                        tocUpdate();
                    }
                })
            }
        })
    }

    disableAllOptions();

    getData(OPTIONS_LIST_FILE).then(res => {
        optionsList = res;
        disabledOptions(optionsList);
        removeActive();
    });

    if (Object.keys(initialData).length) {
        for (const key in initialData) {
            $(`.button__wrap.button--${key}`).removeClass('disabled');
            const buttons = $('.button__wrap').find(`[data-key="${key}"]`);

            buttons.each((_, button) => {
                $(button).removeClass('active');
                if ($(button).attr('data-value') === initialData[key])
                    $(button).addClass('active');
            })
        }
        removeActive();
    }

    function getDataFromHTML() {
        const data = {};
        $('.btn.active').each((_, btn) => {
            if ($(btn).attr('data-key')) {
                const key = $(btn).attr('data-key');
                data[key] = $(btn).attr('data-value');
            }
        })
        return data;
    }

    function disableAllOptions() {
        $('.button__wrap').each((_, item) => {
            $(item).addClass('disabled');
        })
    }

    function removeActive() {
        $('.button__wrap.disabled').find('.btn.btn_o.active').removeClass('active');
    }

    function disabledOptions(options) {
        // FIXME: refactor this function
        if (options === null) {
            return;
        }

        if ($('.button__wrap').length === 0) {
            return;
        }

        const buttons = $('.button__wrap').find(`[data-key="${options['option']}"]`);

        if (typeof options['values'] === 'object' && Object.keys(options['values']).length) {
            $(`.button__wrap.button--${options['option']}`).removeClass('disabled');

            buttons.each((_, item) => {
                $(item).addClass('disabled');
            })

            for (const valuesKey in options['values']) {
                if ($('.button__wrap').find(`[data-key="${options['option']}"][data-value="${valuesKey}"]`)) {
                    $('.button__wrap').find(`[data-key="${options['option']}"][data-value="${valuesKey}"]`).removeClass('disabled');
                }
            }

            if (Object.keys(options['values']).length) {
                if ($('.button__wrap').find('.active.disabled').length) {
                    buttons.filter('.btn.btn_o:not(.disabled)').first().addClass('active');
                    $('.button__wrap').find('.active.disabled').removeClass('active');
                    getDataFromHTML();
                } else if (buttons.filter('.active').length === 0) {
                    buttons.filter('.btn.btn_o:not(.disabled)').first().addClass('active');
                }
            }
        }
        if (typeof options['values'] === 'object' && Object.keys(options['values']).length) {
            const activeOption = buttons.filter('.btn.btn_o.active').attr('data-value');

            disabledOptions(options['values'][activeOption]);
        }
    }

    function setParams() {
        let params = '';

        for (const key in getDataFromHTML()) {
            params += `${key}=${getDataFromHTML()[key]}&`
        }
        return params.slice(0, -1);
    }

    $('.button__wrap').each((_, item) => {
        $(item).click((e) => {
            e.preventDefault();

            if ($(e.target).hasClass('btn')) {
                setActive(e.target, e.currentTarget);

                disableAllOptions();
                disabledOptions(optionsList);
                removeActive();
                pathToData = getUrl(getDataFromHTML(), e.target);
                const url = new URL(window.location);
                window.history.replaceState(null, null, `${url.pathname}?${setParams()}`);

                getData(DATA_FILE).then(res => {
                    if (res[pathToData]) {
                        $.ajax({
                            type: "GET",
                            url: `${res[pathToData]}`,
                            success: function (response) {
                                $('#configurator-content').html(response);
                                tocUpdate();
                            }
                        })
                    }
                })
            }
        })
    })

    // TODO check if below handlers are deprecated (setRepoPath is no defined, ids are unknown).
    const configBtn = $('#config-ci-runners');
    const viewWrap = $('#partial-ci-runners');

    $('#project-config').click((e) => {
        e.preventDefault();

        getData(DATA_FILE).then(res => {
            setRepoPath(res[pathToData]);
        })
    })

    configBtn.click((e) => {
        e.preventDefault();

        let str = '';
        const obj = getDataFromHTML();

        for (const key in obj) {
            str += `${key}-${obj[key]}_`
        }
        str = str.slice(0, -1);

        getData(DATA_FILE).then(res => {
            viewWrap.text(res[str].settings);
        })
    })
})

// URL fetcher with memoization.
const getData = (() => {
    let data_cache = {}
    return async (url) => {
        if (url in data_cache) {
            return data_cache[url];
        }
        const result = await fetch(url)
        data_cache[url] = await result.json()
        return data_cache[url]
    }
})();

function getUrl(obj, target) {
    let str = '';
    const key = $(target).attr('data-key');
    obj[key] = $(target).attr('data-value');

    for (const key in obj) {
        str += `${key}-${obj[key]}_`
    }
    str = str.slice(0, -1);
    return str;
}

function setActive(target, currentTarget) {
    if ($(target).hasClass('btn')) {
        $(currentTarget).find('a.btn').removeClass('active');
        $(target).addClass('active');
    }
}

function tocUpdate() {
    $('#toc').toc({minimumHeaders: 0, listType: 'ul', showSpeed: 0, headers: 'h2, h3:not(.tabs__content>h3), h4'});

    if ($('#toc').is(':empty')) { $('#toc').hide(); }

    /* this offset helps account for the space taken up by the floating toolbar. */
    $('#toc').on('click tap', 'a', function() {
        const target = $(this.getAttribute('href'));
        const scrollTarget = target.offset().top;

        $(window).scrollTop(scrollTarget - 110);
        return false
    })
}

