const slide = document.getElementById("slide");
const prev = document.getElementById("prev");
const next = document.getElementById("next");
const indicator = document.getElementById("indicator");
const lists = document.querySelectorAll(".list");
const totalSlides = lists.length;
let count = 0;
const calcTranslateX = (count) => {
  return `translateX(-${(100 * count) / totalSlides}%)`;
};

slide.style.width = `${100 * totalSlides}%`;

slide
  .querySelectorAll("div")
  .forEach((e) => (e.style.width = `${100 / totalSlides}%`));

const updateListBackground = () => {
  for (let i = 0; i < lists.length; i++) {
    lists[i].style.backgroundColor =
      i === count % totalSlides ? "#000" : "#fff";
  }
};
const nextClick = () => {
  slide.classList.remove(`slide${(count % totalSlides) + 1}`);
  count++;
  if (count >= totalSlides) count = 0;
  slide.classList.add(`slide${(count % totalSlides) + 1}`);
  updateListBackground();
  slide.style.transform = calcTranslateX(count);
};
const prevClick = () => {
  slide.classList.remove(`slide${(count % totalSlides) + 1}`);
  count--;
  if (count < 0) count = totalSlides - 1;
  slide.classList.add(`slide${(count % totalSlides) + 1}`);
  updateListBackground();
  slide.style.transform = calcTranslateX(count);
};
next.addEventListener("click", () => {
  nextClick();
});
prev.addEventListener("click", () => {
  prevClick();
});
indicator.addEventListener("click", (event) => {
  if (event.target.classList.contains("list")) {
    const index = Array.from(lists).indexOf(event.target);
    slide.classList.remove(`slide${(count % totalSlides) + 1}`);
    count = index;
    slide.classList.add(`slide${(count % totalSlides) + 1}`);
    updateListBackground();
    slide.style.transform = calcTranslateX(count);
  }
});

const minimumDistance = 30;
let startX = 0;
let endX = 0;

slide.addEventListener("mousedown", (e) => {
  startX = e.pageX;
});

slide.addEventListener("mouseup", (e) => {
  endX = e.pageX;
  const distanceX = Math.abs(endX - startX);

  if (distanceX < minimumDistance) {
    return;
  }

  if (endX - startX > 0) {
    prevClick();
  } else {
    nextClick();
  }
});

slide.addEventListener("touchstart", (e) => {
  startX = e.touches[0].pageX;
});

slide.addEventListener("touchend", (e) => {
  endX = e.changedTouches[0].pageX;
  const distanceX = Math.abs(endX - startX);

  if (distanceX < minimumDistance) {
    return;
  }

  if (endX - startX > 0) {
    prevClick();
  } else {
    nextClick();
  }
});

const checkDelete = () => {
  const key = "delete";
  return window.prompt(`削除する場合は、${key}と入力してください`) === key;
};
