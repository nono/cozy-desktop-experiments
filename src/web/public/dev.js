;(function (window, document) {
  const $ = document.querySelector.bind(document)
  const $$ = document.querySelectorAll.bind(document)

  // Setting is-active on the active link should be done in the view, but kemal
  // doesn't offer an easy way to do that, so JS, it is.
  $$(".sidebar a").forEach((link) => {
    if (link.href === window.location.href) {
      link.classList.add("is-active")
    }
  })
})(window, document)
