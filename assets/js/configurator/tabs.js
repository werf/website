$(document).ready(() => {

    // Handle tab switches.
    $(document).on('click tap', '.configurator__views-buttons', (e) => {
        e.preventDefault();
        const target = $(e.target);
        setOffsetElements();

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
        setOffsetElements();


        switch (target.attr('class')) {
            case 'file-name':
                setActiveFile(target.attr('data-file-name'), target.parent())
                break;
            case 'file-icon':
                setActiveFile(target.siblings().attr('data-file-name'), target.parent())
                break;
            case 'file__wrap':
                setActiveFile(target.find('.file-name').attr('data-file-name'))
                break;
            case 'folder-name':
            case 'folder-icon':
                toggleHidden(target.parent().parent())
                break;
            case 'folder':
                toggleHidden(target.parent())
                break;
            default:
                setActiveFile(target.find('.file-name').attr('data-file-name'))
        }

        function setActiveFile(dataAttr, parent = null) {
            if (!!structure.find('.file__wrap.active')) {
                structure.find('.file__wrap.active').removeClass('active');
            }
            if (!!filesView.find('.file-view.active')) {
                filesView.find('.file-view.active').removeClass('active');
            }
            if (parent) {
                parent.addClass('active');
            } else {
                target.addClass('active');
            }
            filesView.find(`[data-file-view="${dataAttr}"]`).addClass('active');
        }

        function toggleHidden(parent) {
            parent.toggleClass('open');

            parent.find('.child').each((_, child) => {
                child.classList.toggle('open');
                child.classList.toggle('hidden');
            });
        }
    })
});

function setOffsetElements() {
    const folders = $('.folder__wrap .folder');
    const files = $('.file__wrap');

    folders.each((_, folder) => {
        const padding = $(folder).parent().attr('data-depth') * 15;
        $(folder).css('padding-left', `${padding}px`);
    });

    files.each((_, file) => {
        const padding = $(file).parent().attr('data-depth') * 15;
        $(file).css('padding-left', `${padding + 15}px`);
    })
}