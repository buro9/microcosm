(() => {
	// Relative formatter see: https://blog.webdevsimplified.com/2020-07/relative-time-format/
	const relativeFormatter = new Intl.RelativeTimeFormat(undefined, {
		numeric: 'auto'
	})

	const DIVISIONS = [
		{ amount: 60, name: 'seconds' },
		{ amount: 60, name: 'minutes' },
		{ amount: 24, name: 'hours' },
		{ amount: 7, name: 'days' },
		{ amount: 4.34524, name: 'weeks' },
		{ amount: 12, name: 'months' },
		{ amount: Number.POSITIVE_INFINITY, name: 'years' }
	]

	const formatTimeAgo = (date) => {
		let duration = (date - new Date()) / 1000

		for (let i = 0; i <= DIVISIONS.length; i++) {
			const division = DIVISIONS[i]
			if (Math.abs(duration) < division.amount) {
				return relativeFormatter.format(Math.round(duration), division.name)
			}
			duration /= division.amount
		}
	}

	// precise date formatter, roughly equiv to moment's format('llll');
	const preciseFormatter = new Intl.DateTimeFormat('en', { dateStyle: 'medium', timeStyle: 'medium' });

	const formatPreciseDate = (date) => {
		return preciseFormatter.format(date)
	}

	function updateTimes() {
		const times = document.getElementsByTagName('time');
		[...times].forEach((i) => {
			if (i.classList.contains('plain')) {
				return;
			}

			const dt = i.getAttribute('datetime');
			if (!dt || dt.trim() == '') {
				return
			}

			// displays relative times
			const td = formatTimeAgo(new Date(dt));
			if (i.innerHTML != td) {
				i.innerHTML = td;
			}

			// title for precise date
			const tt = i.getAttribute('title');
			const tip = formatPreciseDate(new Date(dt))
			if (!tt || tt != tip) {
				i.setAttribute('title', tip);
			}
		});
	}

	const TIME_UPDATE_INTERVAL = 60 * 1000; // Update every minute

	const initialiseTimeClickHandler = () => {
		document.body.addEventListener('click', (e) => {
			if (!e.target.matches('time')) {
				return;
			}

			const times = document.getElementsByTagName('time');

			[...times].forEach(i => {
				if (i.classList.contains('plain')) {
					return;
				}

				const ancestor = i.parentNode.parentNode;

				// not sure what this does?
				// probably prevents toggling something in datepicker or similar
				if (ancestor.classList.contains('pills-event') || ancestor.classList.contains('cell-meta-event')) {
					return;
				}

				const title = i.getAttribute('title');
				const html = i.innerHTML;

				i.setAttribute('title', html);
				i.innerHTML = title;
			});

		});
	}

	const makeCodeLookPretty = () => {
		const acceptedLangs = [
			"bsh", "c", "cc", "cpp", "cs", "csh", "cyc", "cv",
			"htm", "html", "java", "js", "m", "mxml", "perl", "pl",
			"pm", "py", "rb", "sh", "xhtml", "xml", "xsl"
		];

		const codeBlocks = document.querySelectorAll('pre > code');

		[...codeBlocks].forEach(codeElem => {
			const pre = codeElem.parentNode;
			codeElem.classList.add('prettyprint');


			pre.classList.add('prettyprint');
			pre.classList.add('linenums');

			// GOTCHA: I don't think the templates even support setting the language of a code block
			const lang = pre.getAttribute('lang');

			if (typeof lang !== 'undefined' && lang != '' && acceptedLangs.indexOf(lang) > -1) {
				pre.classList.add('lang-' + lang);
			}
		});

		prettyPrint();
	}

	const paginationFormsBlurHandler = () => {
		const paginationForms = document.querySelectorAll('form[name=paginationByOffset]');

		[...paginationForms].forEach(f => {
			f.addEventListener('submit', (e) => {
				console.log("Page jump requested")

				const initial = e.target.getAttribute('data-initial'),
					limit = e.target.getAttribute('data-limit'),
					max = e.target.getAttribute('data-max'),
					value = parseInt(e.target.querySelector('input[type=text]').value),
					hidden = e.target.querySelector('input[name=offset]');

				console.log("limit = " + limit + " , max = " + max + " , value = " + value);

				if (!isNaN(value) && value >= 1 && value <= max && value != initial) {
					console.log("Jump")
					if (limit && value) {
						hidden.value = (limit * (value - 1));
					}
				} else {
					console.log("Cancel jump")
					e.preventDefault();
				}
			});

			f.querySelector('input[type=text]').addEventListener('blur', (e) => {
				f.dispatchEvent(new Event('submit', {
					'bubbles': true,
					'cancelable': true,
				}));
			});
		});
	};

	const onDomReadyHandler = () => {
		// times
		updateTimes();
		setInterval(updateTimes, TIME_UPDATE_INTERVAL); // Update every minute
		initialiseTimeClickHandler();

		// code blocks
		makeCodeLookPretty();

		// pagination forms blur handler
		paginationFormsBlurHandler();
	};


	document.addEventListener('DOMContentLoaded', onDomReadyHandler, false);
})();