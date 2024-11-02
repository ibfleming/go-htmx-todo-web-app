const colors = require("tailwindcss/colors");

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./templates/*.templ"],
  plugins: [require("@tailwindcss/forms"), require("@tailwindcss/typography")],
};
