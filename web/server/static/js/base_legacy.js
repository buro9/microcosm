
////////////////////
//	   tooltip    //
////////////////////
// GOTCHA: I think we can remove this and just rely on the title attribute
// (function () {
// 	$('[data-toggle=tooltip]').tooltip({
// 		container: 'body'
// 	});
// })();

////////////////////
//   btn-groups   //
////////////////////
// GOTCHA: This is only relevant for memberships
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