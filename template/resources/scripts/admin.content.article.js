
var adminWebCenter = {}

$(document).ready(function() {

	// Sidebar Accordion Menu:
	$("#main-nav li ul").hide(); // Hide all sub menus
	$("#main-nav li a.current").parent().find("ul").slideToggle("slow"); // Slide down the current

	$("#main-content .main-content-box").hide();
	var currrentDiv = $("#main-nav li a.current").attr('href');
	$(currrentDiv).slideToggle("slow");
	
	// menu item's sub menu
	$("#main-nav li a.nav-top-item").click(	// When a top menu item is clicked...
		function() {
			$(this).parent().siblings().find("a.nav-top-item").removeClass("current");
			$(this).parent().siblings().find("ul").hide();
			$(this).parent().siblings().find("ul").slideUp("normal"); // Slide up all sub menus except the one clicked
			$(this).next().slideToggle("normal"); // Slide down the clicked sub menu
			$(this).addClass("current");
			
			var currentDiv = $(this).attr('href');			
			$(currentDiv).siblings('.main-content-box').hide();
			$(currentDiv).find(".content-box-content div.notification").hide();
			$(currentDiv).show();
			
			return false;
		}
	);

	$("#main-nav li ul li a.nav-sub-item").click(
		// When a top menu item is clicked...
		function() {
			$(this).parent().siblings().find("a.nav-sub-item").removeClass("current"); // Remove class from all tabs
			$(this).addClass("current"); // Add class "current" to clicked tab
			
			var currentDiv = $(this).attr('href');
			$(currentDiv).siblings('.content-box').find("fieldset input.text-input").parent().show();
			$(currentDiv).siblings('.content-box').find("fieldset input.text-input").val("");
			$(currentDiv).siblings('.content-box').find("fieldset textarea.wysiwyg").wysiwyg("setContent", "");
			var options = $(currentDiv).siblings('.content-box').find("fieldset select").get(0);
			if (options) {
				options.selectedIndex = 0;
			} 
			
			$(currentDiv).siblings(".content-box").find(".content-box-header .content-box-tabs").hide();
			$(currentDiv).siblings(".content-box").find(".content-box-content").hide();
			
			$(currentDiv).find(".content-box-header .content-box-tabs a").removeClass("current");
			$(currentDiv).find(".content-box-header .content-box-tabs a.default-tab").addClass("current");
			$(currentDiv).find(".content-box-header .content-box-tabs").show();
			
			var tabContent = $(currentDiv).find(".content-box-header .content-box-tabs a.default-tab").attr("href");
			$(tabContent).siblings(".tab-content").hide();
			
			$(tabContent).show();
			$(currentDiv).find(".content-box-content div.notification").hide();
			$(currentDiv).find(".content-box-content").show();			
			
			return false;
		}
	);

	// Sidebar Accordion Menu Hover Effect:
	$("#main-nav li .nav-top-item").hover(
		function() {
			$(this).stop().animate({paddingRight : "25px"}, 200);
		}, 
		function() {
			$(this).stop().animate({paddingRight : "15px"});
		}
	);

	// Minimize Content Box
	$(".content-box-header h3").css({"cursor" : "s-resize"}); // Give the h3 in Content Box Header a different cursor
	$(".closed-box .content-box-content").hide(); // Hide the content of the header if it has the class "closed"
	$(".closed-box .content-box-tabs").hide(); // Hide the tabs in the header if it has the class "closed"

	/*
	$(".content-box-header h3").click( 
		// When the h3 is clicked...
		function() {
			$(this).parent().next().toggle(); // Toggle the Content Box
			$(this).parent().parent().toggleClass("closed-box"); // Toggle the class "closed-box" on the content box
			$(this).parent().find(".content-box-tabs").toggle(); // Toggle the tabs
		}
	);*/

	// Content box tabs:
	$('.content-box .content-box-content div.tab-content').hide(); // Hide the content divs
	$('ul.content-box-tabs li a.default-tab').addClass('current'); // Add the class "current" to the default tab
	$('.content-box-content div.default-tab').show(); // Show the div with class "default-tab"

	$('.content-box ul.content-box-tabs li a').click(
		// When a tab is clicked...
		function() {
			$(this).parent().siblings().find("a").removeClass('current'); // Remove "current" class from all tabs
			$(this).addClass('current'); // Add class "current" to clicked tab
			
			var currentTab = $(this).attr('href'); // Set variable "currentTab" to the value of href of clicked tab
			$(currentTab).siblings().find("div.notification").hide();
			$(currentTab).siblings().hide(); // Hide all content divs
			
			$(currentTab).show(); // Show the content div with the id equal to the id of clicked tab
			return false;
		}
	);
	
	// Close button:
	$(".close").click(
		function() {
			$(this).parent().fadeTo(400, 0, 
				function() { // Links with the class "close" will close parent
					$(this).slideUp(400);
				}
			);
			
			return false;
		}
	);

	// Alternating table rows:
	$('tbody tr:even').addClass("alt-row"); // Add class "alt-row" to even table rows

	// Check all checkboxes when the one in a table head is checked:
	$('.check-all').click(
		function() {
			$(this).parent().parent().parent().parent().find("input[type='checkbox']").attr('checked',$(this).is(':checked'));
		}
	);

	// Initialise Facebox Modal window:
	$('a[rel*=modal]').facebox(); // Applies modal window to any link with attribute rel="modal"

	// Initialise jQuery WYSIWYG:

	$(".wysiwyg").wysiwyg(); // Applies WYSIWYG editor to any textarea with the class "wysiwyg"
	
	content.initialize(adminWebCenter.accesscode, $("#content"));
	account.initialize(adminWebCenter.accesscode, $("#account"));
});