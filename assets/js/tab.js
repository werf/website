function openTab(evt, linksClass, contentClass, contentId) {
  var i, tabcontent, tablinks;

  tabcontent = document.getElementsByClassName(contentClass);
  for (i = 0; i < tabcontent.length; i++) {
    tabcontent[i].style.display = "none";
  }

  tablinks = document.getElementsByClassName(linksClass);
  for (i = 0; i < tablinks.length; i++) {
    tablinks[i].className = tablinks[i].className.replace(" active", "");
  }

  document.getElementById(contentId).style.display = "block";
  evt.currentTarget.className += " active";
}

// TODO: refactor this code
// $(document).ready(() => {
//   const $buttons = $('[data-tabs-button]');
//   let fileName = '';
//
//   let player = AsciinemaPlayer.create(`/assets/videos/build.cast`, document.getElementById('demo'), {
//     fit: "both",
//     autoplay: true,
//   });
//
//   $buttons.each((idx, button) => {
//     $(button).click(() => {
//       $buttons.each((i, button) => $(button).removeClass('active'))
//       $(button).addClass('active');
//       fileName = $(button).attr('data-tabs-button');
//       player.dispose();
//
//       player = AsciinemaPlayer.create(`/assets/videos/${fileName}.cast`, document.getElementById('demo'), {
//         fit: "both",
//         autoplay: true,
//       });
//     })
//   })
// })
