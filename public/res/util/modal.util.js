// document.addEventListener("DOMContentLoaded", function () {
//   const body = document.body;
//   const modalBtns = document.querySelectorAll(".underlay-btn");

//   modalBtns.forEach(btn => {
//       btn.addEventListener("change", function () {
//           if (btn.checked) {
//               body.style.overflow = "hidden"; // Hide scrollbar
//           } else {
//               setTimeout(() => {
//                   body.style.overflow = "";
//               }, 1500);
//           }
//       });
//   });
// });

// document.addEventListener("DOMContentLoaded", function () {
//   const body = document.body;
//   const modalBtns = document.querySelectorAll(".underlay-btn");

//   function isAnyChecked() {
//       return [...modalBtns].some(btn => btn.checked);
//   }

//   modalBtns.forEach(btn => {
//       btn.addEventListener("change", function () {
//           if (isAnyChecked()) {
//               body.style.overflow = "hidden"; // Hide scrollbar
//           } else {
//               setTimeout(() => {
//                   body.style.overflow = ""; // Restore scrollbar
//               }, 1500);
//           }
//       });
//   });
// });