// -Make entire card clickable
document.addEventListener('DOMContentLoaded', () => {
  document.addEventListener('click', e => {
      const card = e.target.closest('.card');
      if (card && !e.target.closest('a, button') && card.dataset.href) {
          location.href = card.dataset.href;
      }
  });

  document.addEventListener('keydown', e => {
      if (e.key === 'Enter') {
          const card = e.target.closest('.card');
          if (card && card.dataset.href) {
              location.href = card.dataset.href;
          }
      }
  });
});


// -Load cards on-scroll
var blogCardsFetching = false;
var Blogs_Page = 2;
var Blogs_Limit = 6;
var noMoreBlogs = false;

function BlogsCard_fetch() {
  if (noMoreBlogs || blogCardsFetching) { return; }

  blogCardsFetching = true;
  $("#blogCards_loading").show().css("opacity", 1);

  $.post("/blog/card-get", { BlogsPage: Blogs_Page, BlogsLimit: Blogs_Limit }, (response) => {

    response = JSON.parse(response);
    if (response.noMoreBlogs) {
      noMoreBlogs = true;
      $("#blogCards_loading").hide();
    } else {
      $("#BlogCards").append(response.html);
      Blogs_Page++;
    }

    $("#blogCards_loading").css("opacity", 0).hide();
    blogCardsFetching = false;
  });
}

BlogsCard_fetch()
$(window).scroll(function () {
  if ($(window).scrollTop() + $(window).height() > $(document).height() - 1200) {
    if (!blogCardsFetching) {
      BlogsCard_fetch();
    }
  }
});