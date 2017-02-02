(function(w,d,$,undefined){

  var ProfilePicture = (function(){

    var profilePicture = function(opts){
      this.el = false;
      if (typeof opts.el !== 'undefined'){
        this.$el   = $(opts.el);
        this.el    = this.$el[0];
        this.label = this.$el.find('label');
        this.input = this.$el.find('input[type=file]');
      }
      this.bind();
    };

    profilePicture.prototype.activateLabel = function(backgroundImage){
      this.label
        .addClass('active')
        .css('background-image', "url(" + backgroundImage + ")");
    };

    profilePicture.prototype.deactivateLabel = function(){
      this.label
        .removeClass('active')
        .css('background-image', "");
    };

    profilePicture.prototype.update = function(files){

      var file;

      if (files.length < 1){
        this.deactivateLabel();
        return;
      }

      file = files[0];

      if (!file.type.match('image.*')){
        this.deactivateLabel();
        return;
      }
      // adds to <input> element
      this.input[0].files = files;

      var reader = new FileReader();

      reader.onload = $.proxy(function(e){
        this.activateLabel(e.target.result);
      },this);

      reader.readAsDataURL(file);

    };

    profilePicture.prototype.changeHandler = function(e){
      this.update(e.target.files);
    };

    profilePicture.prototype.dragHandler = function(e){
      e.stopPropagation();
      e.preventDefault();
    };

    profilePicture.prototype.dropHandler = function(e){
      e.stopPropagation();
      e.preventDefault();
      this.update(e.originalEvent.dataTransfer.files);
    };

    profilePicture.prototype.bind = function(){

      var events = [
        ['change',  'input[type=file]', 'changeHandler'],
        ['drop',    'label',            'dropHandler'],
        ['dragover','label',            'dragHandler']
      ];

      for(var i in events){
        this.$el.on(events[i][0], events[i][1], $.proxy(this[events[i][2]], this) );
      }
    };

    return profilePicture;

  })();

  if (window.File && window.FileReader && window.FileList && window.Blob){
    window.ProfilePicture = ProfilePicture;
  }

})(window, document, jQuery);
