(function (w, d, undefined) {

  function ArraytoFileList(files) {
    const dt = new DataTransfer();
    files.forEach(function (file) {
      dt.items.add(file);
    })
    return dt.files;
  }

  class FileHandler {

    constructor(opts) {

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

    removeFile(index) {
      const sf = this.stack.slice(index, index + 1)[0];
      this.stack = this.stack.filter(function (i) {
        return i !== sf;
      });

      const files = Array.from(this.input.files).filter(function (f) {
        return !(sf.name === f.name && sf.size === f.size && sf.lastModified === f.lastModified);
      });

      this.input.files = ArraytoFileList(files);

      if (typeof this.onRemove !== 'undefined' && typeof this.onRemove === 'function') {
        this.onRemove(this.stack);
      }
    }

    clear() {
      this.stack = [];
      this.input.files = ArraytoFileList([]);
    }

    parse(filesRaw) {

      var reader, callback;

      const files = filesRaw.filter(function (f) {
        const match = this.stack.find(function (sf) {
          return sf.name === f.name && sf.size === f.size && sf.lastModified === f.lastModified;
        });

        return !match;
      }.bind(this));

      if (files.length < 1) {
        return;
      }

      this.callback_counter = files.length;

      // ugly way of keeping track of the reader.onload async events
      // we only want to call our ondragged callback when all "files" have been loaded
      callback = (function (e, i) {

        const f = files[i];
        const fileProps = {
          lastModified: f.lastModified,
          name: f.name,
          size: f.size,
          type: f.type,
          fileRef: files[i],
        };

        // instance of progressevent assumes readasDataurl was triggered
        if (e instanceof ProgressEvent) {
          const modified_attachment = Object.assign(fileProps, {
            data: e.target.result,
          });
          // we use Array.unshift here to push image files to the front of the stack (ie. opposite of Array.push)
          // this makes it easier when we render to html (ie. will render all images first, then non-images)
          this.stack.unshift(modified_attachment);
        } else {
          this.stack.push(fileProps);
        }

        this.callback_counter--;

        if (this.callback_counter <= 0) {
          if (typeof this.onDragged !== 'undefined' && typeof this.onDragged === 'function') {
            this.onDragged(this.stack);
          }

          this.event_type = false;
          this.input.files = ArraytoFileList(this.stack.map(function (f) {
            return f.fileRef;
          }));
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

    }

    onDragged(fn) {
      if (typeof fn === 'function') {
        this.onDragged = fn;
      }
      return this;
    }

    onRemove(fn) {
      if (typeof fn === 'function') {
        this.onRemove = fn;
      }
      return this;
    }

    clickHandler(e) {
      // noop
    }

    changeHandler(e) {
      if (!this.event_type) {
        this.event_type = "changed";
        this.parse(Array.from(e.target.files));
      }
    }

    dragHandler(e) {
      e.stopPropagation();
      e.preventDefault();
    }

    dropHandler(e) {
      const fileList = e.dataTransfer.files;

      e.stopPropagation();
      e.preventDefault();

      if (!this.event_type) {
        this.event_type = "dropped";
        this.input.files = fileList
        this.parse(Array.from(fileList));
      }
    }

    bind() {

      const events = [
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

      this.el.addEventListener('dragover', e => this.dragHandler(e));
      this.el.addEventListener('drop', e => this.dropHandler(e));

    }


  }

  w.FileHandler = FileHandler;

})(window, document);
