// -Moving nav stick-horizontal
$(document).ready(function () {
  $(window).on('load resize', function () {
    var $thisnav = $('.nav-m .current-menu').offset().left;

    $('.nav-m .menu a').hover(function () {
      var $left = $(this).offset().left - $thisnav;
      var $width = $(this).outerWidth();
      var $start = 0;
      $('.nav-m .initbar').css({ 'left': $left, 'width': $width });
    }, function () {
      var $initwidth = $('.nav-m .current-menu').width();
      $('.nav-m .initbar').css({ 'left': '0', 'width': $initwidth });
    });
  });
});
