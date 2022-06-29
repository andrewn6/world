/** @type {import('tailwindcss').Config} */
module.exports = {
	content: [
		"./pages/**/*.{js,ts,jsx,tsx}",
		"./components/**/*.{js,ts,jsx,tsx}",
    "./animations/**/*.{js,ts,jsx,tsx}",
	],
	theme: {
    fontFamily: {
      sans: ["IBM Plex Sarif", ...fontFamily.sans],
    },
		extend: {},
	},
	plugins: [],
};
