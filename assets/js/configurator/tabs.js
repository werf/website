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

    const structure = document.querySelector('.directory-structure');
    const filesView = document.querySelector('.files-view__wrap');


    structure.addEventListener('click', (e) => {
        if (e.target.classList.contains('file-name')) {
            if (!!structure.querySelector('.file__wrap.active')) {
                structure.querySelector('.file__wrap.active').classList.remove('active');
            }
            if (!!filesView.querySelector('.file-view.active')) {
                filesView.querySelector('.file-view.active').classList.remove('active');
            }
            e.target.parentNode.classList.add('active');
            filesView.querySelector(`.file-view.${e.target.dataset.fileName}`).classList.add('active');
        } else if (e.target.classList.contains('folder-name')) {
            const parent = e.target.parentNode.parentNode;
            parent.querySelectorAll('.child').forEach(child => {
                child.classList.toggle('hidden');
            });
        }
    })
});