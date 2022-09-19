(function (w, d, $, undefined) {
  'use strict';

  w.simpleEditor = (function () {

    var simpleEditor = function (options) {
      this.el = false;
      if (typeof options.el !== "undefined") {
        this.$el = typeof options.el == "string" ? $(options.el) : options.el;
        this.el = this.$el[0];
      }

      this.static_url = $('meta[name="subdomain"]').attr('content');
      if (typeof options.static_url !== 'undefined') {
        this.static_url = options.static_url;
      }

      this.no_attachments = false;
      if (typeof options.no_attachments !== 'undefined') {
        this.no_attachments = options.no_attachments;
      }


      this.textarea = this.$el.find('textarea')[0];

      this.form = this.$el.find('form');

      this.bind();

      return this;
    };


    simpleEditor.prototype.nothingSelected = function () {
      return this.el.selectionStart === this.el.selectionEnd;
    };

    simpleEditor.prototype.getSelectionDetailsObject = function () {


      var text = this.textarea.value,
        startPos = this.textarea.selectionStart,
        endPos = this.textarea.selectionEnd,
        selectedLength = this.textarea.selectionEnd - this.textarea.selectionStart;

      var startText = text.substr(0, startPos),
        selectedText = text.substr(startPos, selectedLength),
        endText = text.substr(endPos, text.length);

      var retval = {
        start: {
          position: startPos,
          text: startText
        },
        end: {
          position: endPos,
          text: endText
        },
        selected: {
          length: selectedLength,
          text: selectedText
        }
      };

      return retval;

    };

    simpleEditor.prototype.applyFormatting = function (text, tag) {

      // splits text into array by newlines and applies tag to each index of array.
      var lines = text.split(/\n/g);
      for (var i = 0; i < lines.length; i++) {
        lines[i] = tag.replace(/%s/g, lines[i]);
      }

      return lines.join('\n');
    };

    simpleEditor.prototype.formattedTextWith = function (tag) {

      var selection = this.getSelectionDetailsObject();

      var newText = selection.start.text +
        this.applyFormatting(selection.selected.text, tag) +
        selection.end.text;

      return newText;
    };

    simpleEditor.prototype.insertLinkWith = function (tag) {

      var selection = this.getSelectionDetailsObject();
      var link, newText;

      if (selection.selected.length < 1) {
        link = w.prompt("Paste url here:");
        if (!link) {
          return false;
        }
      } else {
        link = selection.selected.text;
      }

      newText = selection.start.text +
        this.applyFormatting(link, tag) +
        selection.end.text;

      return newText;
    };

    simpleEditor.prototype.h1 = function () {
      this.textarea.value = this.formattedTextWith("\n%s\n====");
    };

    simpleEditor.prototype.bold = function () {
      this.textarea.value = this.formattedTextWith("**%s**");
    };

    simpleEditor.prototype.italics = function () {
      this.textarea.value = this.formattedTextWith("*%s*");
    };

    simpleEditor.prototype.list = function () {
      this.textarea.value = this.formattedTextWith("*%s");
    };

    simpleEditor.prototype.quote = function () {
      this.textarea.value = this.formattedTextWith("> %s");
    };

    simpleEditor.prototype.quoteAll = function () {
      this.textarea.value = this.applyFormatting(this.textarea.value, "> %s");
    };

    simpleEditor.prototype.link = function () {
      var output = this.insertLinkWith("[%s](%s)");
      if (output) {
        this.textarea.value = output;
      }
    };

    simpleEditor.prototype.image = function () {

      var output = this.insertLinkWith("![%s](%s)");
      if (output) {
        this.textarea.value = output;
      }
    };


    simpleEditor.prototype.fetchAttachments = function (options) {

      var ajaxOptions = {
        url: "",
        type: 'GET'
      };

      if (typeof options !== 'undefined' &&
        typeof options == 'object') {
        ajaxOptions = $.extend({}, ajaxOptions, options);
      }

      return $.ajax(ajaxOptions);

    };

    simpleEditor.prototype.clearAttachmentGallery = function (e) {
      this.$el.find('.reply-box-attachments-gallery').html("");
      this.fileHandler.clear();
    };


    simpleEditor.prototype.renderAttachmentGallery = function (files) {
      var ul, li, img, a, span,
        recognised_file_exts = ['jpg', 'jpeg', 'gif', 'png', 'bmp'],
        gallery = this.$el.find('.reply-box-attachments-gallery');


      ul = document.createElement('ul');
      // check if <ul> already exists if so re-use
      var gallery_ul = gallery.find('ul');
      if (gallery_ul.length > 0) {
        gallery_ul.find('.new-attachment').remove();
        ul = gallery_ul[0];
      }


      /*
      *   note 1: this function checks for two states of data hence the various
      *   checks below. The two states have different data structures.
      *   1) directly from the form (check for file.type)
      *   2) fetched via ajax from /attachments api endpoint (check for file.fileHash)
      *
      *   note 2: data-rel is used for removing/deleting an attachment from the comment
      *   note 3: data-idx is used for removing/deleting a new attachment from the new files stack
      *   note 3: data-new is used to decide whether it's an existing one or not
      */

      if (files.length > 0) {
        for (var i = 0, j = files.length; i < j; i++) {

          li = document.createElement('li');

          if ((typeof files[i].type !== 'undefined' && files[i].type.match('image.*')) ||
            (typeof files[i].fileExt !== 'undefined' && recognised_file_exts.indexOf(files[i].fileExt) > -1)) {
            img = document.createElement('img');

            if (typeof files[i].meta !== 'undefined') {
              if (typeof files[i].meta.links !== 'undefined') {
                img.src = this.static_url + files[i].meta.links[0].href;
              }
            } else {
              img.src = files[i].data;
            }
            if (typeof files[i].fileHash !== 'undefined') {
              img.setAttribute('data-rel', files[i].fileHash);
              li.setAttribute('data-new', false);
            } else {
              img.setAttribute('data-rel', files[i].name || files[i].fileName);
              li.className = li.className + " new-attachment";
              li.setAttribute('data-idx', i);
              li.setAttribute('data-new', true);
            }

            img.name = files[i].name || files[i].fileName;

            li.appendChild(img);

          } else { // non-image attachments
            a = document.createElement('a');

            li.className = "attachments-gallery-row";

            if (typeof files[i].fileHash !== 'undefined') {
              a.setAttribute('data-rel', files[i].fileHash);
              li.setAttribute('data-new', false);
            } else {
              a.setAttribute('data-rel', files[i].name || files[i].fileName);
              li.className = li.className + " new-attachment";
              li.setAttribute('data-idx', i);
              li.setAttribute('data-new', true);
            }
            a.name = files[i].name || files[i].fileName;
            a.innerHTML = files[i].name || files[i].fileName;

            li.appendChild(a);
          }

          span = document.createElement('span');
          span.className = 'remove';
          span.innerHTML = "&times;";
          li.appendChild(span);

          ul.appendChild(li);
        }
      }
      if (gallery_ul.length < 1) {
        this.$el.find('.reply-box-attachments-gallery').html(ul);
      }
    };


    simpleEditor.prototype.removeAttachmentFile = function (e) {
      var self = $(e.currentTarget),
        parent = self.parent(),
        isSavedFile = !parent.data('new'),
        fileToBeRemoved = parent.find('[data-rel]');

      var delete_confirm = window.confirm("Are you sure you want to remove this attachment?");

      if (delete_confirm) {
        if (fileToBeRemoved.length > 0) {
          if (typeof this.attachments_delete == 'undefined') {
            this.attachments_delete = [];
          }

          var field_attachments_delete = this.form.find('input[name="attachments-delete"]');

          if (field_attachments_delete.length < 1) {
            field_attachments_delete = $('<input name="attachments-delete" type="hidden"/>');
            this.form.append(field_attachments_delete);
          }

          var notAlreadyDeleted = this.attachments_delete.indexOf(fileToBeRemoved.attr('data-rel')) === -1;

          if (isSavedFile && notAlreadyDeleted) {
            this.attachments_delete.push(fileToBeRemoved.attr('data-rel'));
            field_attachments_delete.val(this.attachments_delete.join(','));
          } else {
            this.fileHandler.removeFile(parent.data('idx'));
          }
        } else {
          // GOTCHA: this should never happen, I guess?
          this.fileHandler.removeFile(self.index());
        }

        parent.remove();
      }

    };


    simpleEditor.prototype.toggleMarkdownPreview = function (e) {
      var preview_window = this.$el.find('.wmd-preview-wrapper'),
        button_bar = this.$el.find('.wmd-button-bar');

      preview_window.css({
        'height': window.getComputedStyle(this.textarea).height,
        'top': button_bar.height()
      });

      preview_window.toggle();
      if (typeof prettyPrint !== 'undefined') {
        prettyPrint();
      }
    }

    simpleEditor.prototype.handleCtrlEnter = function (e) {
      if ((e.ctrlKey || e.metaKey) && (e.keyCode == 13)) {
        this.form.find('input[type=submit]').click();
      }
    };

    simpleEditor.prototype.bind = function () {

      // create a new instance of markdown.editor.js

      var converter = new Markdown.Converter();
      var help = { handler: function () { alert('help!'); } };
      this.editor = new Markdown.Editor(converter, this.$el, help);
      this.editor.run();
      window.editor = this.editor;


      //////////////////////
      //   attachments    //
      //////////////////////
      if (typeof FileHandler !== 'undefined' && !this.no_attachments) {

        this.fileHandler = new FileHandler({
          el: this.$el.find('.reply-box-attachments')[0],
          dropzone: '.reply-box-drop-zone, .reply-box-attachments-gallery'
        });
        // this callback fires when dropped event fires and all files are collected
        this.fileHandler.onDragged($.proxy(function (files) {
          this.renderAttachmentGallery(files);
        }, this));

        this.fileHandler.onRemove($.proxy(function (files) {
          this.renderAttachmentGallery(files);
        }, this));


        // check for data-num-attachments attribute, if 1 or more
        // fire off an ajax call for the attachments json
        // on a response, pass the attachment items to the renderAttachmentGallery()
        var num_attachments = parseInt(this.$el.attr('data-num-attachments'), 10);

        if (!isNaN(num_attachments) && num_attachments > 0) {
          this.$el.find('.reply-box-attachments-gallery').html("Loading attachments...");

          this.fetchAttachments({
            url: this.$el.attr('data-source') + "attachments"
          }).success($.proxy(function (response) {
            if (typeof response.data.attachments !== 'undefined') {
              this.renderAttachmentGallery(response.data.attachments.items);
            }
          }, this));

        }
      }


      //////////////////////
      //   textcomplete   //
      //////////////////////

      var subdomain = $('meta[name="subdomain"]').attr('content'),
        static_url = subdomain,
        dataSource = subdomain + '/api/v1/profiles?disableBoiler&top=true&q=';

      if (typeof $.fn.textcomplete !== 'undefined') {

        $(this.textarea).textcomplete([
          {
            match: /\B([@+][\-+\w]*)$/,
            search: function (term, callback) {

              var _term = term.substr(1, term.length - 1),
                _symbol = term[0],
                _callback = callback;

              $.ajax({
                url: dataSource + _term,
              }).success(function (data) {

                _callback($.map(data.profiles.items, function (person) {
                  if (person.profileName.toLowerCase().indexOf(_term.toLowerCase()) === 0) {
                    person.symbol = _symbol;

                    return person;
                  } else {
                    return null;
                  }
                }));
              });
            },
            template: function (person) {

              var img_src = static_url + person.avatar;

              if (person.avatar.match(/http[s*]:\/\//)) {
                img_src = person.avatar;
              }
              return '<img src="' + img_src + '" /> ' + person.symbol + person.profileName;
            },
            replace: function (person) {
              return person.symbol + person.profileName;
            },
            index: 1,
            maxCount: 5
          }
        ]);
      }

      ////////////////////////
      //   form validation  //
      ////////////////////////
      if (typeof FormValidator !== 'undefined') {

        if (this.form.length > 0) {
          new FormValidator(
            this.form[0],
            {
              rules: {
                'markdown': ['not_empty', 'maxlength']
              },
              tests: {
                'maxlength': function (field) { var $field = $(field); return $field.val().length < parseInt($field.attr('maxlength'), 10); }
              },
              error_messages: {
                'markdown:not_empty': "* Cannot be empty",
                'markdown:maxlength': "* Cannot exceed max length."
              }
            }
          );
        }
      }


      ////////////////////////
      //     bind events    //
      ////////////////////////
      var events = [
        ['reset', 'form', 'clearAttachmentGallery'],
        ['click', '.reply-box-attachments-gallery li span.remove', 'removeAttachmentFile'],
        ['keydown', 'textarea', 'handleCtrlEnter'],
        ['click', '.wmd-preview-button', 'toggleMarkdownPreview']
      ];

      for (var i in events) {
        this.$el.on(events[i][0], events[i][1], $.proxy(this[events[i][2]], this));
      }
    };

    return simpleEditor;
  })();

})(window, document, jQuery);