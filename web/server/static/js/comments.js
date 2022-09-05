(function(){

  var Comments = (function(){

    var comments = function(options){

      this.$el = false;
      if(typeof options.el !== "undefined"){
        this.$el = $(options.el);
      }

      this.defaultContainer = false;
      if(typeof options.defaultContainer !== "undefined"){
        this.defaultContainer = $(options.defaultContainer);
      }

      this.template = false;
      if(typeof options.template !== "undefined"){
        this.template = options.template;
      }

      this.stack = [];
      this.bind();
    };

    comments.prototype.cleanup = function(e){

      var view;

      for(var i=0,j=this.stack.length;i<j;i++){
        view = this.stack.pop();
        view.off().remove();
      }

      var oldInstances = this.$el.find('.insertReplyBox.active');

      if (oldInstances.length > 0){
        for( i=0,j=oldInstances.length;i<j;i++){

          if (typeof oldInstances[i].$ !== "undefined" &&
              typeof oldInstances[i].$.comment_box !== "undefined" &&
              oldInstances[i].$.comment_box){

            oldInstances[i].$.comment_box.remove();
            oldInstances[i].$.comment_box = false;
            oldInstances.removeClass('active');
          }
        }
      }

    };

    comments.prototype.toggleDefaultContainer = function(){

      if (this.stack.length > 0){
        this.defaultContainer.hide();
      }else{
        this.defaultContainer.show();
      }

    };

    comments.prototype.generateNewInstanceCommentBox = function(options){

      var fragment   = $( this.template ),
          action     = "",
          replyto_id = "",
          id         = "",
          num_attachments = 0;

      if (typeof options.action !== 'undefined'){
        action = options.action;
      }
      fragment.find('form').attr('action',action);


      if (typeof options.id !== 'undefined'){
        id = options.id;
      }
      fragment.find('input[name=id]').val(id);

      if(typeof options.ref !== 'undefined'){
        replyto_id = options.ref;
      }
      fragment.find('input[name=inReplyTo]').val(replyto_id);

      this.cleanup();

      if (typeof options.data_source !== 'undefined'){
        fragment.attr('data-source',options.data_source);
      }
      if (typeof options.num_attachments !== 'undefined'){
        fragment.attr('data-num-attachments',options.num_attachments);
      }

      fragment.simpleEditor = new simpleEditor({
        el : fragment
      });

      if(typeof options.ref !== 'undefined' && options.ref !== ""){
        // auto quote if applicable
        var selectedText = this.getWindowSelectedText();
        if (selectedText){
          fragment.simpleEditor.textarea.value = selectedText;
          fragment.simpleEditor.quoteAll();
        }
      }

      this.stack.push(fragment);

      return fragment;
    };

    comments.prototype.fetchSource = function(url){

      // FIXME: possible bug with relative path, need to get absolute
      return $.ajax({
        url  : url,
        type : 'GET'
      });

    };

    comments.prototype.clickHandler = function(e){

      var _this = e.currentTarget,
          commentBoxOptions;

      if (typeof _this.$ == 'undefined'){
        _this.$  = $(_this);
      }

      commentBoxOptions = {
        action : _this.$.attr('data-action') || "create",
        ref    : _this.$.attr('data-ref') || "",
        id     : _this.$.attr('data-comment-id') || "",
        data_source : _this.$.attr('data-source') || "",
        num_attachments : _this.$.attr('data-num-attachments') || 0
      };

      if (!_this.$.hasClass('active')){

        _this.$.comment_box = $('<div class="generated-comment-box"></div>');
        _this.$.comment_box.append( this.generateNewInstanceCommentBox(commentBoxOptions) );

        _this.$
          .addClass('active')
          .parent().append( _this.$.comment_box );

        // FIXME: not flexible, could be better
        if(_this.$.attr('data-source')){

          _this.$.comment_box.find('input[type=submit]').val('Save changes');

          _this.$.comment_box.find('textarea').attr('placeholder','Loading... Please wait...');

          this.fetchSource(_this.$.attr('data-source')+"source/")
              .success($.proxy(function(response){
                this.$.comment_box
                      .find('textarea')
                      .attr('placeholder','Enter your text here...')
                      .val(response.data.markdown);
              },_this))
              .error($.proxy(function(){
                this.$.comment_box
                      .find('textarea')
                      .attr('placeholder','Enter your text here...');
              },_this));
        }

      }else{
        this.cleanup();
      }
      this.toggleDefaultContainer();

    };

    comments.prototype.getWindowSelectedText = function(){

      var selection     = window.getSelection();

      if (selection && typeof selection.toString == 'function') {
        return selection.toString().replace(/\n/g, "\n> ");
      }

      return false;
    };

    comments.prototype.reset = function(){
      this.cleanup();
      this.toggleDefaultContainer();
    };

    comments.prototype.bind = function(){

      var events = [
        ['click', '.insertReplyBox', 'clickHandler'],
        ['reset', 'form',            'reset']
      ];

      for(var i in events){
        this.$el.on(events[i][0], events[i][1], $.proxy(this[events[i][2]], this) );
      }

    };

    return comments;

  })();

  window.Comments = Comments;

})();