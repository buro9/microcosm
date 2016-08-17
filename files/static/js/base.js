function updateTimes() {
	$('time').each(function(i) {
		var t = $(this);

		if (t.hasClass("plain")) {return;}

		var dt = t.attr('datetime');
		if (dt && dt.trim() != '') {
			var m = moment(dt);
			var f = t.attr('format');
			// Relative date for simplicity
			var td = (f && f.trim() != '') ? m.format(f) : m.fromNow();
			if (t.text() != td) {
				t.text(td);
			}
			// Tooltip for precise date
			var tt = t.attr('title');
			var tip = m.format('llll');
			if (!tt || tt != tip) {
				t.attr('title', tip);
			}
		}
	});
}

$('document').ready(function() {
	// Updates <time> to use relative times
	updateTimes();
	setInterval(updateTimes, 60000); // Update every minute

	// Make code look pretty
	$('pre > code').addClass('prettyprint')
			.parent().addClass('prettyprint').addClass('linenums');

	var acceptedLangs = ["bsh", "c", "cc", "cpp", "cs", "csh", "cyc", "cv", "htm", "html", "java", "js", "m", "mxml", "perl", "pl", "pm", "py", "rb", "sh", "xhtml", "xml", "xsl"];

	$('pre.prettyprint').each(function(index) {
		var lang = $(this).attr('lang');

		if (typeof lang !== 'undefined' && lang != '' && $.inArray(lang, acceptedLangs) > -1) {
			$(this).addClass('lang-' + lang);
		}
	});

	prettyPrint();

	// toggle <time> html -> title -> html
	$('body').on('click', 'time', function() {
		if ($(this).parent().parent().hasClass('pills-event') || $(this).parent().parent().hasClass('cell-meta-event')) {
			return;
		}

		$('time').each(function(ii) {
			var t = $(this);
			if (t.hasClass("plain")) {return;}

			title = t.attr('title');
			html  = t.html();
			t.html(title).attr('title',html);
		});
	});
});

////////////////////
//	pagination    //
////////////////////
(function(){
	$('form[name=paginationByOffset]').submit(function(e){
		console.log("Page jump requested")

		var self    = $(this),
			initial = self.attr('data-initial'),
			limit   = self.attr('data-limit'),
			max     = self.attr('data-max'),
			value   = parseInt(self.find('input[type=text]').val()),
			hidden  = self.find('input[name=offset]');

		console.log("limit = " + limit + " , max = " + max + " , value = " + value);

		if (!isNaN(value) && value >= 1 && value <= max && value != initial) {
			console.log("Jump")
			if (limit && value){
				hidden.val(limit * (value-1));
			}
		} else {
			console.log("Cancel jump")
			e.preventDefault();
		}
	});
	$('form[name=paginationByOffset] > input[type=text]').blur(function() {
		$(this).parent().submit();
	});
})();

////////////////////
//	   tooltip    //
////////////////////
(function(){
	$('[data-toggle=tooltip]').tooltip({
		container : 'body'
	});
})();

////////////////////
//   btn-groups   //
////////////////////
(function(){

	var btn_groups = '.btn-group';

	var toggleButtonParent = function(e){
		var self 				= $(e.currentTarget),
				siblings 		= $( 'input[name="' + self.attr('name') + '"]' ),
				activeClass = 'active';

		if (self.is(':checked')){
			siblings.parent().removeClass( activeClass );
			self.parent().addClass( activeClass );
		}else{
			self.parent().removeClass( activeClass );
		}

	}

	$(btn_groups).on('change','input[type=radio]', toggleButtonParent);

})();