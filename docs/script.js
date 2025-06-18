function showLanding() {
  document.getElementById("landing-page").style.display = "block";
  document.getElementById("help-page").style.display = "none";
}

function showHelp() {
  document.getElementById("landing-page").style.display = "none";
  document.getElementById("help-page").style.display = "block";
}

function showPopup(platform) {
  alert(`${platform} version is not yet available for download.`);
}

// Help section navigation
function showHelpSection(sectionId) {
  const sections = document.querySelectorAll(".help-section");
  sections.forEach((section) => section.classList.remove("active"));

  const buttons = document.querySelectorAll(".help-nav-btn");
  buttons.forEach((btn) => btn.classList.remove("active"));

  document.getElementById(sectionId).classList.add("active");

  event.target.classList.add("active");
}

// FAQ toggle
function toggleFAQ(element) {
  const answer = element.nextElementSibling;
  const icon = element.querySelector("span");

  if (answer.classList.contains("active")) {
    answer.classList.remove("active");
    icon.textContent = "+";
  } else {
    answer.classList.add("active");
    icon.textContent = "-";
  }
}

document.querySelectorAll('a[href^="#"]').forEach((anchor) => {
  anchor.addEventListener("click", function (e) {
    e.preventDefault();
    const target = document.querySelector(this.getAttribute("href"));
    if (target) {
      target.scrollIntoView({
        behavior: "smooth",
      });
    }
  });
});
