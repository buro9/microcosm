/////////
// MAP //
/////////
(function(){

  var EventFormMap = (function(){

    var eventFormMap = function(opts){

      // reference to map marker
      this.marker = [];

      // container element
      this.el = null;
      if (typeof opts.el !== 'undefined'){
        if (typeof opts.el == 'string'){
          this.el  = opts.el;
          this.$el = $(this.el);
        }else{
          this.$el = opts.el;
          this.el = this.$el[0];
        }
      }

      // maps cordinates to <input>s
      this.mappings = {
        'lat'   : { el : 'input[name=lat]',   value : null },
        'lng'   : { el : 'input[name=lon]',   value : null },
        'north' : { el : 'input[name=north]', value : null },
        'south' : { el : 'input[name=east]',  value : null },
        'east'  : { el : 'input[name=south]', value : null },
        'west'  : { el : 'input[name=west]',  value : null }
      };

      if (typeof opts.mappings !== 'undefined'){
        this.mappings = $.extend({}, this.mappings, opts.mappings);
      }

      this.map = L.map( 'map-canvas' , { zoomAnimation : false })
                  .setMaxBounds([[-90,-180],[90,180]]); // Restrict map to valid lat:lng pairs

      var googleLayer = new L.Google('ROADMAP');
      this.map.addLayer(googleLayer);

      var osmCycleLayer = new L.TileLayer("https://a.tile.thunderforest.com/cycle/{z}/{x}/{y}.png");
      var osmLandscapeLayer = new L.TileLayer("https://a.tile.thunderforest.com/landscape/{z}/{x}/{y}.png");

      this.map.addControl(new L.Control.Layers({
        'Open Cycle Map' : osmCycleLayer,
        'Landscape'      : osmLandscapeLayer,
        'Google Maps'    : googleLayer
        })
      );

      this.bind();

      return this;

    };

    eventFormMap.prototype.clearMarkers = function(){
      if (this.marker.length > 0){
        for(var i=0,j=this.marker.length;i<j;i++){
          this.map.removeLayer(this.marker.pop());
        }
      }
      return this;
    };

    eventFormMap.prototype.addMarker = function(latlng){

      var newMarker = L.marker(latlng);
      this.marker.push(newMarker);

      return this;
    };

    eventFormMap.prototype.renderMarkers = function(){
      if (this.marker.length > 0){
        for(var i=0,j=this.marker.length;i<j;i++){
          this.map.addLayer( this.marker[i] );
        }
      }
      return this;
    };

    eventFormMap.prototype.dropMarker = function(e){

      this.clearMarkers()
          .addMarker(e.latlng)
          .renderMarkers()
          .update(e.latlng)
          .save();

      return this;
    };

    // If this is the edit screen and we need to set the map to known values then
    // we do so here
    // @param lat number
    // @param lng number
    // @param bounds array(array(north, west),array(south, east))
    eventFormMap.prototype.restore = function(lat, lng, bounds){

      this.map.fitBounds(bounds);
      // Bounds will always be 1 level too far out due to imprecision in the
      // numbers and Leaflet aggressively ensuring the bounds fit inside the
      // map area. To prevent us slowly zooming out with each edit, we zoom in
      this.map.zoomIn();

      this.clearMarkers();
      this.addMarker(new L.LatLng(lat, lng));
      this.renderMarkers();

      //assumes form inputs will be prefilled
    };


    eventFormMap.prototype.updateMapBounds = function(){
      // Get the bounds (map box)
      var b  = this.map.getBounds(),
          sw = b.getSouthWest(),
          ne = b.getNorthEast();

      this.mappings['north'].value  = ne.lat;
      this.mappings['south'].value  = ne.lng;
      this.mappings['east'].value   = sw.lat;
      this.mappings['west'].value   = sw.lng;

      return this;
    };
    eventFormMap.prototype.updateLatlng = function(latlng){

      this.mappings['lat'].value    = latlng.lat;
      this.mappings['lng'].value    = latlng.lng;

      return this;
    };

    eventFormMap.prototype.update = function(latlng){

      this.updateLatlng(latlng);
      this.updateMapBounds();

      return this;
    };

    eventFormMap.prototype.save = function(){
      var input;
      for(var i in this.mappings){
        if (typeof this.mappings[i].value !== 'undefined'){
          input = $(this.mappings[i].el);
          if (input.length > 0){
            input.val(this.mappings[i].value);
          }
        }
      }
      return this;
    };

    eventFormMap.prototype.onMapDragHandler = function(e){
      this.updateMapBounds().save();
    };


    // resets location inputs, clears markers
    eventFormMap.prototype.reset = function(){

      var input;

      for(var i in this.mappings){
        input = $(this.mappings[i].el);
        if (input.length > 0){
          input.val('');
        }
        this.mappings[i].value = null;
      }

      this.clearMarkers();

      return this;

    };

    eventFormMap.prototype.bind = function(){

      var events = [];

      if (events.length>0){
        for(var i=0,j=events.length;i<j;i++){
          this.$el.on(events[i][0], events[i][1], $.proxy(this[events[i][2]], this) );
        }
      }

      // for map object
      this.map.on('click', $.proxy(this.dropMarker,this) );
      this.map.on('dragend', $.proxy(this.onMapDragHandler, this));

    };


    return eventFormMap;

  })();

  window.EventFormMap = EventFormMap;



})();


/////////////////////////////////////////
//   event calendar and time controls  //
/////////////////////////////////////////
(function(){
  'use strict';

  var EventDateForm = (function(){

    var EVENT_DATE_TYPE_TBA      = 0,
        EVENT_DATE_TYPE_SINGLE   = 1,
        EVENT_DATE_TYPE_MULTIPLE = 2;

    var eventDateForm = function(opts){

      this.el = false;

      if(typeof opts.el !== "undefined"){
        this.$el = $(opts.el);
        this.el  = this.$el[0];
      }

      // controls
      this.controls = {
        start_calendar            : null,
        start_calendar_start_time : null,
        start_calendar_end_time   : null,
        end_calendar              : null,
        end_calendar_end_time     : null
      };

      var i=0;

      if (typeof opts.controls !== "undefined"){
        for (i in opts.controls){
          this.controls[i] = $(opts.controls[i]);
        }
      }

      this.form = {
        when      : null,
        duration  : null
      };

      if (typeof opts.form !== "undefined"){
        for (i in opts.form){
          this.form[i] = $(opts.form[i]);
        }
      }

      // defaults
      this.startDate = new Date();
      this.endDate   = new Date();

      this.event_date_type = EVENT_DATE_TYPE_SINGLE;

      if (typeof opts.startDate !== "undefined"){
        // Force UTC values to be displayed by negating current offset as we
        // received a UTC value
        var _date = new Date(opts.startDate);
        var userOffset = _date.getTimezoneOffset()*60000;
        this.startDate = new Date(_date.getTime()+userOffset);
      }

      // with duration, we just work out the end date
      // and show "multiple" state if endDate is not on the same day as startDate
      if (typeof opts.duration !== "undefined"){

        this.endDate = new Date(this.startDate.getTime() + parseInt(opts.duration,10) * 60000);

        if( this.endDate.getDate()+""+this.endDate.getMonth()+""+this.endDate.getYear() !==
            this.startDate.getDate()+""+this.startDate.getMonth()+""+this.startDate.getYear()
        ){
          this.event_date_type = EVENT_DATE_TYPE_MULTIPLE;
        }
      }

      this.updateStartTimesUI(getTimeStringFromDate(this.startDate));
      this.updateEndTimesUI(getTimeStringFromDate(this.endDate));

      this.controls.start_calendar.date = new Date(this.startDate.getTime());
      this.controls.end_calendar.date = new Date(this.endDate.getTime());

      this.updateDatePickerUI(this.controls.start_calendar);
      this.updateDatePickerUI(this.controls.end_calendar);

      this.toggleCalendars();

      this.bind();
    };

    // Add a time that is within a string in the format "03:45 PM" onto an existing
    // date object
    function setTime(d, ts) {
      if (!d) {return;}
      var time = ts.match(/(\d?\d):(\d\d)\s(P?)/);
      d.setHours(parseInt(time[1],10) + ((parseInt(time[1],10) < 12 && time[3]) ? 12 : 0));
      d.setMinutes(parseInt(time[2],10) || 0);
    }

    // Given a date object, return the time formatted as "03:45 PM"
    function getTimeStringFromDate(d) {
      var hh = d.getHours();
      var mi = d.getMinutes();
      var pm = (hh > 11);

      // round minutes to nearest 15min block
      if (mi > 45){
        hh = hh+1;
        mi = 0;
        pm = (hh > 11);
      } else if (mi > 30) {
        mi = 45;
      } else if (mi > 15) {
        mi = 30;
      } else if (mi > 0) {
        mi = 15;
      } else {
        mi = 0;
      }

      return '' + (hh <= 9 ? '0' + hh : (hh > 12 ? hh - 12: hh)) + ':' +
        (mi <= 9 ? '0' + mi : mi) + ' ' +
        (pm ? 'PM' : 'AM');
    }
    eventDateForm.prototype.getTimeStringFromDate = getTimeStringFromDate;

    // template object for calendar display
    // @param  dateObject - a javascript Date object
    // @return string     - string which will be converted and injected into the dom
    function template(dateObject){

      var d = dateObject,
          locale_days   = ['Sunday','Monday','Tuesday','Wednesday','Thursday','Friday','Saturday'],
          locale_months = [
            "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"
          ],
          output = "";

      output = '<div class="item-event-date">' +
                 '<span class="item-event-date-digit">' + d.getDate() + '</span>' +
                 '<span class="item-event-date-day">' + locale_days[d.getDay()] +
                   '<strong>'+ locale_months[d.getMonth()] + ' ' + d.getFullYear() + '</strong>' +
                 '</span>' +
               '</div>';

      return output;
    }

    // Not all browsers know how to do toISOString() as it was only added in
    // ECMAScript 5, so we just ship our own version.
    function dateToUtcString(d) {
      function pad(n) { return n < 10 ? '0' + n : n; }
      return d.getUTCFullYear()           + '-' +
             pad(d.getUTCMonth()    + 1)  + '-' +
             pad(d.getUTCDate())    + 'T' +
             pad(d.getUTCHours())   + ':' +
             pad(d.getUTCMinutes()) + ':' +
             pad(d.getUTCSeconds()) + 'Z';
    }

    eventDateForm.prototype.setEventType = function(type){

      this.event_date_type = type;

    };

    eventDateForm.prototype.getStartTime = function(){
      return this.controls.start_calendar_start_time.val().trim();
    };

    eventDateForm.prototype.getEndTime = function(){

      var endTimeControl;

      if (this.event_date_type == EVENT_DATE_TYPE_SINGLE){
        endTimeControl = this.controls.start_calendar_end_time;
      }else if (this.event_date_type == EVENT_DATE_TYPE_MULTIPLE){
        endTimeControl = this.controls.end_calendar_end_time;
      }else{
        console.log('getEndDateTime(): unrecognized event_date_type');
      }
      return endTimeControl.val().trim();
    };

    eventDateForm.prototype.calculateEventDuration = function(){

      var duration,
          startDateTime,
          endDateTime;

      // We can glue them together
      startDateTime = this.startDate;
      setTime(startDateTime, this.getStartTime());

      if (this.event_date_type == EVENT_DATE_TYPE_SINGLE){
        endDateTime = new Date(this.startDate.getTime());
      }else if (this.event_date_type == EVENT_DATE_TYPE_MULTIPLE){
        endDateTime = this.endDate;
      }else{
        console.log('calculateEventDuration(): unrecognized event_date_type');
      }
      setTime(endDateTime, this.getEndTime());

      // To calculate the difference between the dateTimes as an integer
      // representing minutes
      duration = (endDateTime.valueOf() / 60000) - (startDateTime.valueOf() / 60000);

      // And if we've got junk input, default to an hour
      if (!duration || duration < 0) {
        duration = 60;
      }

      return {
        startDateTime : startDateTime,
        endDateTime   : endDateTime,
        duration      : parseInt(duration,0)
      };

    };

    eventDateForm.prototype.updateEventFormWhenDuration = function(){

      var eventDate,
          eventDateFields;

      eventDate = this.calculateEventDuration();

      // Force saving of UTC by removing local offset
      var _date = new Date(eventDate.startDateTime);
      var userOffset = _date.getTimezoneOffset()*60000;
      var startDate = new Date(_date.getTime()-userOffset);

      eventDateFields = [
        [ this.form.when,     dateToUtcString(startDate) ],
        [ this.form.duration, eventDate.duration ]
      ];

      for(var i=0,j=eventDateFields.length;i<j;i++){
        eventDateFields[i][0].val(eventDateFields[i][1]);
      }
    };

    eventDateForm.prototype.resetEventFormWhenDuration = function(){

      var eventDateFields;

      eventDateFields = [
        this.form.when,
        this.form.duration
      ];

      for(var i=0,j=eventDateFields.length;i<j;i++){
        eventDateFields[i].val('');
      }
    };


    // events

    // updates event dates for post request
    // @param ev   - expects a specific event object returned by bootstrap-datepicker plugin
    //             - at a minimum { currentTarget: <domelement>, date: <jsDateObject> }
    eventDateForm.prototype.updateEventStartDate = function(ev){

      this.startDate = ev.date;
      this.controls.start_calendar.date = new Date(this.startDate.getTime());

      this.updateEventFormWhenDuration();
      this.updateDatePickerUI(this.controls.start_calendar);
    };

    eventDateForm.prototype.updateEventEndDate = function(ev){

      this.endDate = ev.date;
      this.controls.end_calendar.date = new Date(this.endDate.getTime());

      this.updateEventFormWhenDuration();
      this.updateDatePickerUI(this.controls.end_calendar);
    };

    // handles the UI of the calendar
    // @param ev   - expects a specific event object returned by bootstrap-datepicker plugin
    //             - at a minimum { currentTarget: <domelement>, date: <jsDateObject> }
    eventDateForm.prototype.updateDatePickerUI = function(calendar){
      calendar.datepicker('hide');
      calendar.html(template(calendar.date));
    };

    eventDateForm.prototype.updateStartTimesUI = function(desiredValue){
      this.controls.start_calendar_start_time.val(desiredValue);
    };

    eventDateForm.prototype.updateEndTimesUI = function(desiredValue){
      this.controls.start_calendar_end_time.val(desiredValue);
      this.controls.end_calendar_end_time.val(desiredValue);
    };


    // if..else based on event_date_type
    // _____________   _____________
    // |      x     |  |      z     |
    // |____________|  |____________|
    //
    // [ x ] to [ y ]  [ z ]
    eventDateForm.prototype.toggleCalendars = function(){

          // cache references to the calendar elements we are interested in
      var elements_event_calendars         = this.$el.find('.form-datepicker'),
          elements_event_date_type_single  = this.$el.find('.form-datepicker-single'),
          elements_event_date_type_mutiple = this.$el.find('.form-datepicker-multiple');

      // shows all "x,y", hides all "z"
      if (this.event_date_type == EVENT_DATE_TYPE_SINGLE){
        elements_event_calendars.show();
        elements_event_date_type_single.show();
        elements_event_date_type_mutiple.hide();

        this.updateEventFormWhenDuration();
      }
      // shows all "x,z", hides all "y"
      if (this.event_date_type == EVENT_DATE_TYPE_MULTIPLE){
        elements_event_calendars.show();
        elements_event_date_type_single.hide();
        elements_event_date_type_mutiple.show();

        this.updateEventFormWhenDuration();

        var mutipleStateRadioButton = this.$el.find('#event-date-type-options input[type=radio]').eq(1);
        if (mutipleStateRadioButton.length>0){
          mutipleStateRadioButton[0].checked = "checked";
        }
      }
      // hides "x,y,z"
      if (this.event_date_type == EVENT_DATE_TYPE_TBA){
        elements_event_calendars.hide();

        this.resetEventFormWhenDuration();
      }
    };

    eventDateForm.prototype.onEventDateTypeToggle = function(e){
      var self = $(e.currentTarget);
      this.setEventType(self.val());
      this.toggleCalendars();
    };

    eventDateForm.prototype.bind = function(){
      var events = [
        [ 'change', '#event-date-type-options input[type=radio]', 'onEventDateTypeToggle' ],
        [ 'change', '.event-time',                                'updateEventFormWhenDuration' ]
      ];

      if (events.length>0){
        for(var i=0,j=events.length;i<j;i++){
          this.$el.on(events[i][0], events[i][1], $.proxy(this[events[i][2]], this) );
        }
      }

      this.controls.start_calendar
        .datepicker()
        .on('changeDate', $.proxy(this.updateEventStartDate,this));

      this.controls.end_calendar
        .datepicker()
        .on('changeDate',   $.proxy(this.updateEventEndDate,this));
    };

    return eventDateForm;

  })();

  window.EventDateForm = EventDateForm;

})();



/*
*   form
*   attendee limit toggle
*/
(function(){
  'use strict';


  var AttendeesForm = (function(){

    var attendeesForm = function(options){

      this.el = null;
      if(typeof options.el !== 'undefined'){
        this.$el = $(options.el);
      }

      this.controls = {
        choices : null
      };
      if(typeof options.choices !== 'undefined'){
        this.controls.choices = $(options.choices);
      }

      this.form = {
        attendees : null
      };
      if(typeof options.attendees){
        this.form.attendees = $(options.attendees);
      }

      this.has_attendees_limit = false;
      if(this.form.attendees.val() !== "" || this.form.attendees.val() !== 0 ){
        this.has_attendees_limit = true;
      }

      this.bind();

    };


    attendeesForm.prototype.enabledAttendeeLimitField = function(){
      this.form.attendees.attr('disabled',false).focus();
    };
    attendeesForm.prototype.disableAttendeeLimitField = function(){
      this.form.attendees.val(0).attr('disabled',true);
    };

    attendeesForm.prototype.onChangeChoiceHandler = function(e){

      var radio = $(e.currentTarget);

      if (radio.val() == "1"){
        this.has_attendees_limit = 1;
        this.enabledAttendeeLimitField();
      }else{
        this.has_attendees_limit = 0;
        this.disableAttendeeLimitField();
      }

    };

    attendeesForm.prototype.bind = function(){

      var events = [
        ['change', 'input[name='+this.controls.choices[0].name+']', 'onChangeChoiceHandler' ]
      ];

      if (events.length>0){
        for(var i=0,j=events.length;i<j;i++){
          this.$el.on(events[i][0], events[i][1], $.proxy(this[events[i][2]], this) );
        }
      }

    };

    return attendeesForm;

  })();

  window.AttendeesForm = AttendeesForm;

})();