(function(){

	// usage
	//
	// var a = new Validator( <html#form>, {
	//   rules : { '<form#field.name>' : 'test_name' },
	//   tests : { 'test_name' : <js#function($jquery#field)> }
	// });
	// 2 builtin tests called 'not_empty' and 'not_duplicate'
	// tests should return false if invalid
	//
	// eg.
	// var validateObject = new Validator(
	//   document.getElementById('myform'),
	//   {
	//     rules : {
	//       'first_name' : 'not_empty',
	//       'last_name'  : ['not_empty', 'minLength3' ]
	//     },
	//     tests : {
	//       'minLength3' : function(field){ return field.val().length > 3; }
	//     }
	//   });
	//
	// on failure, preventsDefault() the form, returns array of errors

	var Validator = (function(){

		var validator = function(form, options){

			this.form   = form;
			this.$form  = $(this.form);

			this.errors = [];
			this.initialState = this.getSerializedForm();

			// tests
			this.tests = {
				'not_empty'     : this.not_empty
			};
			if (typeof options.tests !== 'undefined'){
				this.tests = $.extend({},this.tests,options.tests);
			}

			// rules
			var user_rules = {};
			if (typeof options.rules !== 'undefined'){
				user_rules = options.rules;
			}
			this.rules = user_rules;

			// messages
			var error_messages = {};
			if (typeof options.error_messages !== 'undefined'){
				error_messages = options.error_messages;
			}
			this.error_messages = error_messages;

			this.bind();

			return this;
		};

		// tests
		validator.prototype.not_empty = function(field){
			return field.val().trim() !== '';
		};

		validator.prototype.getSerializedForm = function(){

			var serializedArray = this.$form.serializeArray();
			serializedArray.shift();// remove the csrf token

			var serializedArrayToString = [];
			for(var i=0,j=serializedArray.length;i<j;i++){
				serializedArrayToString.push( serializedArray[i].name + '=' + serializedArray[i].value );
			}
			return $.trim(serializedArrayToString.join('&'));
		};

		validator.prototype.not_duplicate = function(value){

			// this function does not diff for input[type=file], so if found, assume true
			var file_inputs = this.$form.find('input[type=file]');

			if (file_inputs.length > 0){
				return false;
			}else{
				return this.getSerializedForm() === this.initialState;
			}
		};


		validator.prototype.validate = function(e){

			this.errors = [];

			// dupe check
			if (this.not_duplicate()){
				this.errors.push('dupe');
			}

			var field;
			for(var field_name in this.rules){
				field = this.$form.find('[name='+field_name+']');

				if (typeof this.rules[field_name] !== 'undefined'){

					if (typeof this.rules[field_name] == 'object'){ // is array

						for(var i=0,j=this.rules[field_name].length;i<j;i++){
							if( !this.tests[this.rules[field_name][i]](field) ){
								this.errors.push([field_name, this.rules[field_name][i]].join(":"));
							}
						}

					}else{ // not an array
						if( !this.tests[this.rules[field_name]](field) ){
							this.errors.push([field_name, this.rules[field_name]].join(":"));
						}
					}
				}

			}

			if (this.errors.length>0){
				e.preventDefault();
				console.log(this.errors);
				this.applyErrorsToForm(this.errors);
				return this.errors;
			}
		};

		validator.prototype.applyErrorsToForm = function(errors){

			var field;

			if (errors.length > 0){
				for(var i=0,j=errors.length;i<j;i++){

					field = errors[i].split(':');
					if (field.length > 1){
						this.addError(field[0],errors[i]);
						this.addErrorStyles(field[0]);
					}
				}
			}

		};

		validator.prototype.addError = function(field_name, error_name){

			var field_input = this.$form.find('[name="'+field_name+'"]'),
					parent      = field_input.parent(),
					error_label;

			if (parent.find('.control-label').length < 1){
				error_label = document.createElement('label');
				error_label.className    = "control-label";
				error_label.textContent = this.error_messages[error_name];
				parent.prepend(error_label);
			}

		};

		validator.prototype.removeError = function(field_name){
			var field_input = this.$form.find('[name="'+field_name+'"]'),
					parent      = field_input.parent(),
					error_label = parent.find('.control-label');

			if (error_label.length>0){
				error_label.remove();
			}

		};

		validator.prototype.addErrorStyles = function(field_name){
			var self   = this.$form.find('[name="'+field_name+'"]'),
					parent = self.parent();
			if (!parent.hasClass('has-error')){
				parent.addClass('has-error');
			}

		};

		validator.prototype.removeErrorStyles = function(field_name){

			var self   = this.$form.find('[name="'+field_name+'"]'),
					parent = self.parent();

			if (parent.hasClass('has-error')){
				parent.removeClass('has-error');
			}

		};

		validator.prototype.onChangeHandler = function(e){

			var self       = $(e.currentTarget),
					field_name = self.attr('name');

			this.removeError(field_name)
			this.removeErrorStyles(field_name);

		};

		validator.prototype.bind = function(){

			this.$form.on('submit', $.proxy(this.validate,this));
			this.$form.on('change', '.has-error input, .has-error textarea', $.proxy(this.onChangeHandler,this));

			this.$form.on('keyup', '.has-error textarea', $.proxy(this.onChangeHandler,this));

		}

		return validator;

	})();

	window.FormValidator = Validator;
})();