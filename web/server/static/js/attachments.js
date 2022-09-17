(function (w, d) {

  var FileHandler = (function () {

    var fileHandler = function (opts) {

      if (typeof opts.el !== 'undefined') {
        if (typeof opts.el === 'string') {
          this.el_name = opts.el;
          this.el = d.querySelector(opts.el);
        } else if (typeof opts.el === 'object') {

          this.el = opts.el;
          this.el_name = '.' + this.el.className;

        } else {
          return false;
        }
      }

      if (typeof opts.dropzone !== 'undefined') {
        this.dropzone = opts.dropzone;
      }
      this.input = this.el.querySelector('input[type=file]');

      this.event_type = false;

      this.stack = [];
      this.bind();

      return this;
    };

    fileHandler.prototype.removeFile = function (index) {
      this.stack.splice(index, 1);
      this.input.files = this.stack;

      if (typeof this.onRemove !== 'undefined' && typeof this.onRemove === 'function') {
        this.onRemove(this.stack);
      }
    };

    fileHandler.prototype.clear = function () {
      for (var i = 0, j = this.stack.length; i < j; i++) {
        this.stack.pop();
      }
      this.input.files = this.stack;
    };

    fileHandler.prototype.parse = function (files) {

      var reader, callback;

      if (files.length < 1) {
        return;
      }

      this.input.files = files;
      this.callback_counter = this.input.files.length;

      // ugly way of keeping track of the reader.onload async events
      // we only want to call our ondragged callback when all "files" have been loaded
      callback = (function (e, i) {

        // instance of progressevent assumes readasDataurl was triggered
        if (e instanceof ProgressEvent) {
          const f = this.input.files[i];

          const modified_attachment = {
            lastModified: f.lastModified,
            name: f.name,
            size: f.size,
            type: f.type,
            data: e.target.result,
            originalFile: f,
          };
          // we use Array.unshift here to push image files to the front of the stack (ie. opposite of Array.push)
          // this makes it easier when we render to html (ie. will render all images first, then non-images)
          this.stack.unshift(modified_attachment);
        } else {
          const modified_attachment = this.input.files[i];
          this.stack.push(modified_attachment);
        }

        this.callback_counter--;
        if (this.callback_counter <= 0) {
          if (typeof this.onDragged !== 'undefined' && typeof this.onDragged === 'function') {
            this.onDragged(this.stack);
          }
          this.event_type = false;
        }
      }).bind(this);


      // loops and reads through all files recieved, skips files which are not images
      for (var i = 0, j = files.length; i < j; i++) {
        if (files[i].type.match('image.*')) {

          reader = new FileReader();

          // we want to call the callback but keep this context of "this" within this
          // object and pass through "i" as a counter
          reader.onload = (function (i) {
            return function (e) {
              callback(e, i);
            };
          })(i);

          reader.readAsDataURL(files[i]);
        } else {
          callback(files[i], i);
        }
      }

      return this;

    };

    fileHandler.prototype.onDragged = function (fn) {
      if (typeof fn === 'function') {
        this.onDragged = fn;
      }
      return this;
    };

    fileHandler.prototype.onRemove = function (fn) {
      if (typeof fn === 'function') {
        this.onRemove = fn;
      }
      return this;
    };

    fileHandler.prototype.clickHandler = function (e) {
      this.input.value = null;
      this.stack = [];
    };

    fileHandler.prototype.changeHandler = function (e) {
      if (!this.event_type) {
        this.event_type = "changed";
        this.parse(e.target.files);
      }
    };

    fileHandler.prototype.dragHandler = function (e) {
      e.stopPropagation();
      e.preventDefault();
    };

    fileHandler.prototype.dropHandler = function (e) {
      e.stopPropagation();
      e.preventDefault();
      if (!this.event_type) {
        this.event_type = "dropped";
        this.parse(e.originalEvent.dataTransfer.files);
      }
    };

    fileHandler.prototype.bind = function () {
      var events = [
        ['change', 'input[type=file]', 'changeHandler'],
        ['click', 'input[type=file]', 'clickHandler']
        // ['drop',      this.dropzone,      'dropHandler'],
        // ['dragover',  this.dropzone,      'dragHandler']
      ];

      for (const event of events) {
        const [eventType, srcSelector, handlerFunc] = event;

        this.el.addEventListener(eventType, (e) => {
          if (!e.target.matches(srcSelector)) {
            return false;
          }

          this[handlerFunc].call(this, e);
        });
      }

      this.el.addEventListener('dragover', this.dragHandler.bind(this));
      this.el.addEventListener('drop', this.dropHandler.bind(this));

    };

    return fileHandler;

  })();

  w.FileHandler = FileHandler;

})(window, document);