$(document).ready(() => {
    const headerHeight = $('.header').height();
    const bannerHeight = $('.guides-banner').outerHeight();
    const breadcrumbs = $('.breadcrumbs-container');
    const breadcrumbsHeight = breadcrumbs.height();
    const fullBreadcrumbsHeight = breadcrumbs.outerHeight(true);
    const breadcrumbsMarginTop = parseInt(breadcrumbs.css('margin-top'));
    const sidebarWrapperInner = $('.sidebar__wrapper-inner');
    const sidebar = $('.sidebar__container');
    let sidebarOffsetTop = 0
    if (sidebar.length > 0) {
        sidebarOffsetTop = sidebar.offset().top - bannerHeight - breadcrumbsHeight + breadcrumbsMarginTop;
    }
    const footerHeight = $('.footer').outerHeight(true);
    const docHeight = $(document).height();
    const screenHeight = $(window).outerHeight();
    let bottomFixPoint = docHeight - (footerHeight + screenHeight);

    if ($(window).scrollTop() > breadcrumbsHeight + breadcrumbsMarginTop + bannerHeight) {
        sidebarWrapperInner.css({
            position: 'fixed',
            top: `${headerHeight + breadcrumbsMarginTop + bannerHeight}px`
        });
    } else {
        setTopOffset($(window).scrollTop(), sidebarOffsetTop, sidebarWrapperInner, headerHeight, breadcrumbsHeight, breadcrumbsMarginTop, fullBreadcrumbsHeight, bannerHeight);
    }

    setFooterOffset($(window).scrollTop(), bottomFixPoint, sidebarWrapperInner, screenHeight, footerHeight, docHeight);

    $(window).scroll(function() {
        const scrolled = $(this).scrollTop();
        bottomFixPoint = $(document).height() - (footerHeight + screenHeight);

        setTopOffset(scrolled, sidebarOffsetTop, sidebarWrapperInner, headerHeight, breadcrumbsHeight, breadcrumbsMarginTop, fullBreadcrumbsHeight, bannerHeight);

        setFooterOffset(scrolled, bottomFixPoint, sidebarWrapperInner, screenHeight, footerHeight, docHeight)
    });
});


function setTopOffset(scrolled, offsetTop, sidebarWrapper, headerHeight, breadcrumbsHeight, breadcrumbsMarginTop, fullBreadcrumbsHeight, bannerHeight) {
    if (scrolled > offsetTop && scrolled < breadcrumbsHeight + breadcrumbsMarginTop + bannerHeight) {
        sidebarWrapper.css({
            position: 'fixed',
            top: `${headerHeight + fullBreadcrumbsHeight + bannerHeight - scrolled}px`
        });
    } else if (scrolled > breadcrumbsHeight + breadcrumbsMarginTop + bannerHeight) {
        sidebarWrapper.css({
            position: 'fixed',
            top: `${headerHeight + breadcrumbsMarginTop}px`
        });
    } else if (scrolled < offsetTop && scrolled < breadcrumbsHeight + breadcrumbsMarginTop + bannerHeight) {
        sidebarWrapper.css({
            top: `${headerHeight + fullBreadcrumbsHeight + bannerHeight - scrolled}px`,
        });
    }
}

function setFooterOffset(scrolled, bottomFixPoint, sidebarWrapper, screenHeight, footerHeight, docHeight) {
    if (scrolled > bottomFixPoint) {
        sidebarWrapper.css({
            bottom: `${scrolled + screenHeight + footerHeight + 25 - docHeight}px`
        })
    } else if (scrolled < bottomFixPoint) {
        sidebarWrapper.css({
            bottom: `25px`
        })
    }
}