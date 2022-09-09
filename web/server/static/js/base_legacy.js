
////////////////////
//	   tooltip    //
////////////////////
(function () {
	$('[data-toggle=tooltip]').tooltip({
		container: 'body'
	});
})();

////////////////////
//   btn-groups   //
////////////////////
(function () {

	var btn_groups = '.btn-group';

	var toggleButtonParent = function (e) {
		var self = $(e.currentTarget),
			siblings = $('input[name="' + self.attr('name') + '"]'),
			activeClass = 'active';

		if (self.is(':checked')) {
			siblings.parent().removeClass(activeClass);
			self.parent().addClass(activeClass);
		} else {
			self.parent().removeClass(activeClass);
		}

	}

	$(btn_groups).on('change', 'input[type=radio]', toggleButtonParent);

})();