(function(){

  var PeopleWidget = (function(){

    var peopleWidget = function(options){

      if (typeof options.el !== 'undefined'){
        this.$el = $(options.el);
      }

      this.people = [];
      if(typeof options.people !== 'undefined'){
        this.people = options.people;
      }
      this.people_invited = [];
      if(typeof options.invited !== 'undefined'){
        this.people_invited = options.invited;
      }

      this.dataSource = null;
      if (typeof options.dataSource !== 'undefined'){
        this.dataSource = options.dataSource;
      }

      this.static_url = "";
      if (typeof options.static_url !== 'undefined'){
        this.static_url = options.static_url;
      }

      this.show_on_focus = true;
      if (typeof options.show_on_focus !== 'undefined'){
        this.show_on_focus = options.show_on_focus;
      }

      this.on = true; // flag to "pause" the widget


      this.isActive = false; // flag to keep the popover from closing

      this.container    = this.createWidgetContainer();
      this.widget_form  = this.createWidgetForm();
      this.widget_list  = this.createWidgetList();

      this.is_textbox = false;
      if (typeof options.is_textbox !== 'undefined' &&
          options.is_textbox === true){
        this.is_textbox = true;
        this.widget_input = this.$el[0];
      }else{
        this.widget_input = this.createWidgetInput();
        this.container.appendChild(this.widget_input);
      }



      this.widget_form.appendChild(this.widget_list);
      this.container.appendChild(this.widget_form);

      if (typeof options.list_target == "undefined"){
        this.popover = this.createPopover();
        this.popover.appendChild(this.container);
        document.body.appendChild(this.popover);
      }else{
        this.popover = false;
        this.list_target = $(options.list_target);
        this.list_target.append(this.container);
      }


      // used for cursor select
      this.cursor = 0;

      this.bind();
    };


    peopleWidget.prototype.createWidgetContainer = function(){
      var widget_container = document.createElement('div');
      widget_container.className = "people-widget";

      return widget_container;
    };

    peopleWidget.prototype.createWidgetForm = function(){
      var widget_form = document.createElement('form');
      return widget_form;
    };

    peopleWidget.prototype.createWidgetInput = function(){
      var widget_input  = document.createElement('input');
      widget_input.type = "text";
      widget_input.placeholder = "Enter a username...";
      widget_input.autocomplete = "off";
      return widget_input;
    };

    peopleWidget.prototype.createWidgetList = function(){
      var widget_list  = document.createElement('ul');
      widget_list.className = "people-widget-list";
      return widget_list;
    };

    peopleWidget.prototype.createSubmitButton = function(){
      var button = document.createElement('input');
      button.type = 'submit';
      button.value = 'Invite';

      return button;
    };

    peopleWidget.prototype.createPopover = function(){

      var popover = document.createElement('div');
      popover.className      = "people-widget-popover";
      popover.style.display  = 'none';
      popover.style.position = 'absolute';

      if (!this.is_textbox){
        nib = document.createElement('span');
        nib.className = 'sprite sprite-nib-up';
        popover.appendChild(nib);
      }
      return popover;
    };

    peopleWidget.prototype.ItemPerson = function(descriptor){
      var link    = document.createElement('a');
      var avatar  = document.createElement('img');

      link.id          = descriptor.id;
      if (descriptor.avatar.indexOf('http') == 0) {
        avatar.src     = descriptor.avatar;
      } else {
        avatar.src     = this.static_url + descriptor.avatar;
      }

      link.appendChild(avatar);
      link.appendChild(document.createTextNode(descriptor.profileName));

      return link;
    };

    peopleWidget.prototype.addPersonToInvitedById = function(id){

      var query_id, invited;

      query_id = parseInt(id,10);
      invited = this.people.filter(function(person){
        return person.id === query_id;
      });

      if (invited.length>0){
        this.people_invited.push(invited[0]);
        this.render();
        if(typeof this.onSelection !== 'undefined' &&
          typeof this.onSelection === 'function'){
          this.onSelection(this.people_invited, invited[0]);
        }
        this.widget_input.value = "";
      }
      this.cursor = -1;
      this.renderCursor();
      return this;
    };

    peopleWidget.prototype.removePersonFromInvitedById = function(id){

      var index_to_remove = false;
      if (this.people_invited.length>0){
        for(var i=0,j=this.people_invited.length;i<j;i++){
          if(this.people_invited[i].id === id){
            index_to_remove = i;
          }
        }
        if (index_to_remove !== false){
          this.people_invited.splice(index_to_remove,1);
        }
      }
      this.cursor = -1;
      this.renderCursor();
      return this;
    };

    peopleWidget.prototype.excludeInvitedPeople = function(list){

      var invited_ids = [],
          new_list    = [];

      if (this.people_invited.length > 0){
        for(var i=0,j=this.people_invited.length;i<j;i++){
          if (typeof this.people_invited !== 'undefined'){
            invited_ids.push(this.people_invited[i].id);
          }
        }
        new_list = list.filter(function(person){
          return invited_ids.indexOf(person.id) === -1;
        });
      }else{
        new_list = list;
      }

      return new_list;
    };


    peopleWidget.prototype.sortPeopleListByName = function(list){
      return list.sort(function(a,b){ return b.name < a.name; });
    };

    peopleWidget.prototype.clearPeopleList = function(){
      this.widget_list.innerHTML = "";
      return this;
    };

    peopleWidget.prototype.renderPeopleList = function(list, options){

      var li, entry;

      for(var i=0,j=list.length;i<j;i++){

        li = document.createElement('li');
        if (typeof options !== 'undefined' &&
            typeof options.className !== 'undefined'){
          li.className = options.className;
        }

        entry = this.ItemPerson(list[i]);
        li.appendChild(entry);
        if (typeof entry.id !== 'undefined'){
          li.rel = entry.id;
        }
        this.widget_list.appendChild(li);
      }
      return this;
    };

    peopleWidget.prototype.render = function(){

      var query, list, nib;

      this.clearPeopleList();

      /* renders the invited list */
      // if (this.people_invited.length > 0){
      //   list = this.sortPeopleListByName(this.people_invited);
      //   this.renderPeopleList(list,{ className : 'invited'});
      // }

      query = this.widget_input.value;

      if ($.trim(query)!==""){

        list = this.excludeInvitedPeople(this.people);

        if (list.length>0){
          this.renderPeopleList(this.sortPeopleListByName(list));

          if (this.is_textbox){
            nib = document.createElement('span');
            nib.className = 'sprite sprite-nib-up';
            this.widget_list.appendChild(nib);
          }
        }else{
          var empty = document.createElement('li');
          empty.textContent = "No results";
          this.widget_list.appendChild(empty);
        }

      }
    };

    peopleWidget.prototype.renderCursor = function(){

      var rows = $(this.widget_list).find('li');

      if (rows.length>0){
        rows.removeClass('active');
        $(rows[this.cursor]).addClass('active');
      }
      return this;
    }

    peopleWidget.prototype.reset = function(){
      this.widget_input.value = "";
      this.clearPeopleList();
      this.people_invited = [];
    };

    peopleWidget.prototype.calcTriggerOffset = function(){

      var offset = this.$el.offset();
      offset.top = offset.top + this.$el.outerHeight();
      offset.left = offset.left + (this.$el.outerWidth()/2);

      return offset;

    };

    peopleWidget.prototype.calcPopoverPosition = function(){

      var offset;

      offset = this.calcTriggerOffset();
      if (this.is_textbox){
        offset.left = offset.left - (this.$el.outerWidth()/2);
      }else{
        offset.top = offset.top + 10;
        offset.left = offset.left - (this.$el.outerWidth()/2);//- (this.popover.offsetWidth/2);
      }

      this.popover.style.top  = offset.top + 'px';
      this.popover.style.left = offset.left + 'px';

      return this;
    };

    peopleWidget.prototype.show = function(){
      if (this.popover){
        this.popover.style.display = "block";
        this.calcPopoverPosition();
      }else{
        if (typeof this.list_target !== 'undefined'){
          this.list_target.show();
        }
      }
      this.widget_input.focus();
    };

    peopleWidget.prototype.hide = function(e){
      if (this.popover){
        this.popover.style.display = "none";
      }else{
        if (typeof this.list_target !== 'undefined'){
          this.list_target.hide();
        }
      }
    };

    peopleWidget.prototype.queryDataSource = function(query, success, error){

      var ajaxOptions;

      if (typeof this.dataSource !== 'undefined'){
        ajaxOptions = {
          url     : this.dataSource + query,
          type    : 'GET',
          success : success,
          error   : error
        };

        $.ajax(ajaxOptions);
      }

      return this;
    };

    peopleWidget.prototype.parseAPIResponse = function(data){

      if (typeof data.profiles !== 'undefined'){
        this.people = data.profiles.items;
      }
      return this;
    };

    peopleWidget.prototype.invitedListToDelimitedString = function(){
      var result;
      result = $.map(this.people_invited,function(person){ return person.id; });
      result = result.join(',');
      return result;
    };

    peopleWidget.prototype.changeHandler = function(e){

      if (this.widget_input.type !== 'textarea'){
        e.preventDefault();
        e.stopPropagation();
      }
      if (!this.on){
        return;
      }

      var currentList = this.excludeInvitedPeople(this.people);


      if ([13].indexOf(e.which)!== -1){

        if (e.which === 13){
          if (currentList.length > 0 && this.cursor > -1){
            this.addPersonToInvitedById(currentList[this.cursor].id);
          }
          this.show();
        }

      }else if ( [38, 40].indexOf(e.which)!==-1 ){


        switch (e.which){
          case 38:
            this.cursor = this.cursor - 1;
            break;
          case 40:
            this.cursor = this.cursor + 1;
            break;
          default:
            break;
        }

        if (this.cursor < 0){
          this.cursor = 0;
        }
        if (this.cursor > currentList.length-1 ){
          this.cursor = currentList.length-1;
        }

        this.renderCursor();

      }else{
        if ($.trim(this.widget_input.value) !== ''){

          var new_query = this.widget_input.value;

          if (this.widget_input.type === 'textarea'){
            var queryRefs = this.widget_input.value.match(/[+@](\w+)\s*?/ig);
            if (queryRefs !== null){
              new_query = queryRefs[queryRefs.length-1].replace(/[+@]/ig,'');
            }
          }

          this.queryDataSource(
            new_query,
            $.proxy(function(data){
              this.parseAPIResponse(data);
              this.render();
              this.renderCursor();
            },this),
            function(e){
              console.log(e);
            }
          );
        }else{
          this.cursor = 0;
          this.people = [];
          this.render();
        }
      }

    };

    peopleWidget.prototype.clickHandler = function(e){
      var self = e.currentTarget;
      e.stopPropagation();
      if (typeof self.rel !== 'undefined'){
        this.addPersonToInvitedById(self.rel);
      }
    };

    peopleWidget.prototype.submitHandler = function(e){
      e.preventDefault();
    };

    peopleWidget.prototype.onSelection = function(fn){
      if (typeof fn === 'function'){
        this.onSelection = $.proxy(fn,this);
      }
      return this;
    };

    peopleWidget.prototype.bind = function(){

      $(this.widget_input)
        .on('keyup',    $.proxy(this.changeHandler,this))
        .on('keypress', $.proxy(function(e){

          if (this.widget_input.type !== 'textarea'){
            if (e.which === 13){
              e.preventDefault();
            }
          }

          if (this.on){
            if (e.which === 13){
              e.preventDefault();
            }
          }

        },this));


      $(this.container).on('click', 'li:not(.invited)', $.proxy(this.clickHandler,this));


      if (this.is_textbox){

        if (this.show_on_focus){
          this.$el.on('focus', $.proxy(this.show,this));
        }

      }else{
        this.$el.on('click', $.proxy(function(){
          if (this.popover.style.display === "block"){
            this.hide();
          }else{
            this.show();
          }
        },this));
      }

      $('body').on('click', $.proxy(function(){ this.hide(); },this));
      $(this.container).on('click', function(e){ e.stopPropagation(); });
      this.$el.on('click', function(e){ e.stopPropagation();});

    };

    return peopleWidget;

  })();

  window.PeopleWidget = PeopleWidget;

})();

//////////////////////
//   Participants   //
//////////////////////

// used to render attendee/invitees lists

(function(){

  var Participating = (function(){

    var participating = function(options){
      if (typeof options !== 'undefined'){
        this.$el = $(options.el);
      }
      if(typeof options.static_url !== 'undefined'){
        this.static_url = options.static_url;
      }
      if (typeof options.className !== 'undefined'){
        this.className = options.className;
      }
    };

    participating.prototype.show = function(){
      this.$el.show();
      return this;
    };

    participating.prototype.hide = function(){
      this.$el.hide();
      return this;
    };

    participating.prototype.render = function(peopleList){

      var ul, li;

      if (peopleList.length > 0){
        ul = document.createElement('ul');
        if (this.className){
          ul.className = this.className;
        }

        for(var i=0,j=peopleList.length;i<j;i++){
          li = document.createElement('li');
          li.rel = peopleList[i].id;
          li.appendChild( PeopleWidget.prototype.ItemPerson.apply(this,[peopleList[i]]) );

          if (typeof peopleList[i].sticky === 'undefined'){
            span = document.createElement('span');
            span.className = "remove";
            span.innerHTML = "&times;";
            span.rel       = peopleList[i].id;
            li.children[0].appendChild(span);
          }

          ul.appendChild(li);
        }

        this.$el.html(ul);
      }

      return this;
    };

    return participating;

  })();

  window.Participating = Participating;

})();