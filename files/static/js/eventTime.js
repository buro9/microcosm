// This file is used to display event times correctly. It is not used on the
// create or edit pages

////////////////
// LOCAL TIME //
////////////////
(function(){
  // A simple case of pulling the UTC values in the page, and allowing the
  // JavaScript locale to dictate the values.

  var startTime = null;
  var endTime = null;
  
  $("time[itemprop='startDate']").each(function() {
    startTime = new Date($(this).attr('datetime'));
    var p = $(this).closest('.item-event-date');
    
    // Set the date to the current locale
    p.find('.item-event-date-digit').first().text(startTime.getDate());
    var days = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
    var months = ['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December'];
    p.find('.item-event-date-day').first().html(days[startTime.getDay()] + '<strong>' + months[startTime.getMonth()] + ' ' + startTime.getFullYear() + '</strong>');

    // Set the time to the current locale
    var t = p.find("time[itemprop='duration']");
    var dur = t.attr('datetime').substring(1,t.attr('datetime').length-1);
    endTime = new Date(startTime.getTime() + dur * 60000);

    if(t.text().indexOf('-') != -1){
      // Show start and end time
      t.text(
        (startTime.getHours() > 12 ? startTime.getHours()-12 : startTime.getHours()) +
        ':' + (startTime.getMinutes() < 10 ? '0' : '') + startTime.getMinutes() +
        ' ' + (startTime.getHours() >= 12 ? 'PM' : 'AM') + 
        ' - ' + 
        (endTime.getHours() > 12 ? endTime.getHours()-12 : endTime.getHours()) +
        ':' + (endTime.getMinutes() < 10 ? '0' : '') + endTime.getMinutes() +
        ' ' + (endTime.getHours() >= 12 ? 'PM' : 'AM')
      );
    } else {
      // Only show start time
      t.text(
        (startTime.getHours() > 12 ? d.getHours()-12 : startTime.getHours()) +
        ':' + (startTime.getMinutes() < 10 ? '0' : '') + startTime.getMinutes() +
        ' ' + (startTime.getHours() >= 12 ? 'PM' : 'AM')
      );
    }
  });

  $("time[itemprop='endDate']").each(function() {
    var p = $(this).closest('.item-event-date');
    
    // Set the date to the current locale
    p.find('.item-event-date-digit').first().text(endTime.getDate());
    var days = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
    var months = ['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December'];
    p.find('.item-event-date-day').first().html(days[endTime.getDay()] + '<strong>' + months[endTime.getMonth()] + ' ' + endTime.getFullYear() + '</strong>');

    // Set the time to the current locale
    p.find(".item-event-date-time").first().text(
      (endTime.getHours() > 12 ? endTime.getHours()-12 : endTime.getHours()) +
      ':' + (endTime.getMinutes() < 10 ? '0' : '') + endTime.getMinutes() +
      ' ' + (endTime.getHours() >= 12 ? 'PM' : 'AM')
    );
  });
})();