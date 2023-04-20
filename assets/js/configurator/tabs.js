$(document).ready(() => {

    // Handle tab switches.
    $(document).on('click tap', '.configurator__views-buttons', (e) => {
        e.preventDefault();
        const target = $(e.target);

        // Handle only clicks on non-active buttons.
        if (target.hasClass('btn') && !target.hasClass('active')) {
            // Activate clicked button.
            target.parent().find('[data-view-button]').removeClass('active');
            target.addClass('active');
            // Activate tab content.
            $('#view__wrap').find('[data-view-content]').removeClass('active');
            const btnAttr = target.attr('data-view-button');
            $(`[data-view-content = ${btnAttr}]`).addClass('active');
        }
    });

    // Toggling details in tabs.
    $(document).on('click tap', '.details__summary', (e) => {
        e.preventDefault();
        $(e.target).closest('.details').toggleClass('active');
    });

});