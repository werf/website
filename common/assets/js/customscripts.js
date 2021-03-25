$('#mysidebar').height($(".nav").height());

$( document ).ready(function() {
    var wh = $(window).height();
    var sh = $("#mysidebar").height();

    if (sh + 100 > wh) {
        $( "#mysidebar" ).parent().addClass("layout-sidebar__sidebar_a");
    }

});

$(document).ready(function () {
  $('.expand-content').hide();

  $('.expand .expand-click').click(function(event){
    event.preventDefault();
    console.log($(this).parent());
    $(this).parent().find('.expand-content').toggle();
  });
});



// Update GitHub stats
$(document).ready(function () {
  var github_requests = [],
  github_stats = JSON.parse(localStorage.getItem('werf_github_stats')) || null;

  function getGithubReuests() {
    $('[data-roadmap-step]').each(function () {
      var $step = $(this);
      github_requests.push($.get('https://api.github.com/repos/werf/werf/issues/' + $step.data('roadmap-step'), function (data) {
        github_stats['issues'][$step.data('roadmap-step')] = (data.state === 'closed');
      }));
    });
    github_requests.push($.get("https://api.github.com/repos/werf/werf", function (data) {
      github_stats['stargazers'] = data.stargazers_count
    }));
    return github_requests;
  }

  function updateGithubStats() {
    $('.gh_counter').each(function () {
      $(this).text(github_stats['stargazers']);
    });
    $('[data-roadmap-step]').each(function () {
      var $step = $(this);
      if (github_stats['issues'][$step.data('roadmap-step')] == true) $step.addClass('roadmap__steps-list-item_closed');
    });
  }

  if (github_stats == null || Date.now() > (github_stats['updated_on'] + 1000 * 60 * 60)) {
    github_stats = {'updated_on': Date.now(), 'issues': {}, 'stargazers': 0};
    $.when.apply($, getGithubReuests()).done(function() {
      updateGithubStats();
      localStorage.setItem('werf_github_stats', JSON.stringify(github_stats));
    });
  } else {
    updateGithubStats();
  }
});

// Guides clipbord copy
var action_toast_timeout;
function showActionToast(text) {
  clearTimeout(action_toast_timeout);
  var action_toast = $('.action-toast');
  action_toast.text(text).fadeIn()
  action_toast_timeout = setTimeout(function(){ action_toast.fadeOut() }, 5000);
}

$(document).ready(function(){
  new ClipboardJS('[data-snippetcut-btn-name]', {
    text: function(trigger) {
      showActionToast('Имя скопировано')
      return $(trigger).closest('[data-snippetcut]').find('[data-snippetcut-name]').text();
    }
  });
  new ClipboardJS('[data-snippetcut-btn-text]', {
    text: function(trigger) {
      showActionToast('Текст скопирован')
      return $(trigger).closest('[data-snippetcut]').find('[data-snippetcut-text]').text();
    }
  });
});

// Guide sticky btns
$(document).ready(function () {
  var $sidebar_btns = $('.sidebar__btns');
  if ($sidebar_btns.length > 0) {
    var sidebar_btns_watcher = scrollMonitor.create($sidebar_btns[0], 125);
    sidebar_btns_watcher.stateChange(function() {
			$sidebar_btns.toggleClass('sticky', sidebar_btns_watcher.isAboveViewport);
		});
    setTimeout(function() {
      $sidebar_btns.toggleClass('sticky', sidebar_btns_watcher.isAboveViewport);
    }, 100); // prevents weird behaiviour when page is refreshed when buttons should be sticky
  }
});

$(document).ready(function () {
  $('.expand_columns_content').hide();
  $('.expand_columns_button').click(function(event){
    event.preventDefault();
    $('.expand_columns_content').hide();
    let name = '#' + this.getAttribute('id') + '__content';
    $(name).show();
  });
});

$(document).ready(function () {
  $('.header__menu-item').removeClass('header__menu-item_active');
  $('.header__menu-item > a').each(function () {
    if (this.href.endsWith("/guides.html")) {
      $(this).parent().addClass('header__menu-item_active')  ;
    }
  });
});
