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

    $(document).on('click tap', '.directory-structure', (e) => {
        const structure = $('.directory-structure');
        const filesView = $('.files-view__wrap');
        const target = $(e.target);


        if (target.hasClass('file-name')) {
            if (!!structure.find('.file__wrap.active')) {
                structure.find('.file__wrap.active').removeClass('active');
            }
            if (!!filesView.find('.file-view.active')) {
                filesView.find('.file-view.active').removeClass('active');
            }
            target.parent().addClass('active');
            filesView.find(`[data-file-view="${target.attr('data-file-name')}"]`).addClass('active');
        } else if (target.hasClass('folder-name')) {
            const parent = target.parent().parent();
            parent.find('.child').each((_, child) => {
                child.classList.toggle('hidden');
            });
        }
    })
});