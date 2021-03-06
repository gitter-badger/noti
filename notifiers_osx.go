// +build darwin

package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa

#import <Foundation/Foundation.h>
#import <objc/runtime.h>
#include <AppKit/AppKit.h>
#include <errno.h>

@implementation NSBundle(noti)
- (NSString *)notiIdentifier {
	return @"com.apple.terminal";
}
@end

@interface NotiDelegate : NSObject<NSUserNotificationCenterDelegate>
@property (nonatomic, assign) BOOL delivered;
@end

@implementation NotiDelegate
- (void) userNotificationCenter:(NSUserNotificationCenter *)center didActivateNotification:(NSUserNotification *)notification {
	self.delivered = YES;
}
- (void) userNotificationCenter:(NSUserNotificationCenter *)center didDeliverNotification:(NSUserNotification *)notification {
	[NSApp activateIgnoringOtherApps:YES];
	self.delivered = YES;
}
@end

void BannerNotify(const char* title, const char* message, const char* sound) {
	errno = 0;
	@autoreleasepool {
		Class nsBundle = objc_getClass("NSBundle");
		method_exchangeImplementations(
			class_getInstanceMethod(nsBundle, @selector(bundleIdentifier)),
			class_getInstanceMethod(nsBundle, @selector(notiIdentifier))
		);

		NotiDelegate *notiDel = [[NotiDelegate alloc]init];
		notiDel.delivered = NO;

		NSUserNotificationCenter *nc = [NSUserNotificationCenter defaultUserNotificationCenter];
		nc.delegate = notiDel;

		NSUserNotification *nt = [[NSUserNotification alloc] init];
		nt.title = [NSString stringWithUTF8String:title];
		nt.informativeText = [NSString stringWithUTF8String:message];

		if ([[NSString stringWithUTF8String:sound] isEqualToString: @"_mute"] == NO) {
			nt.soundName = [NSString stringWithUTF8String:sound];
		}

		[nc deliverNotification:nt];

		while (notiDel.delivered == NO) {
			[[NSRunLoop currentRunLoop] runUntilDate:[NSDate dateWithTimeIntervalSinceNow:0.1]];
		}
	}
}
*/
import "C"

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"unsafe"
)

const (
	specificPart = `
    NOTI_SOUND
        Banner success sound. Default is Ping. Possible options are Basso, Blow,
        Bottle, Frog, Funk, Glass, Hero, Morse, Ping, Pop, Purr, Sosumi,
        Submarine, Tink. See /System/Library/Sounds for available sounds.
    NOTI_SOUND_FAIL
        Banner failure sound. Default is Basso. Possible options are Basso,
        Blow, Bottle, Frog, Funk, Glass, Hero, Morse, Ping, Pop, Purr, Sosumi,
        Submarine, Tink. See /System/Library/Sounds for available sounds.
    NOTI_VOICE
        Name of voice used for speech notifications. See "say -v ?" for
        available voices.

BUGS
    Banner notifications don't fire in tmux.

    Clicking on banner notifications causes unexpected behavior.`
)

func init() {
	flag.Usage = func() {
		fmt.Printf(manual, specificPart)
	}
}

// bannerNotify triggers an NSUserNotification.
func bannerNotify(n notification) error {
	var sound string
	if n.failure {
		sound = os.Getenv(soundFailEnv)
		if sound == "" {
			sound = "Basso"
		}
	} else {
		sound = os.Getenv(soundEnv)
		if sound == "" {
			sound = "Ping"
		}
	}

	t := C.CString(n.title)
	m := C.CString(n.message)
	s := C.CString(sound)
	defer C.free(unsafe.Pointer(t))
	defer C.free(unsafe.Pointer(m))
	defer C.free(unsafe.Pointer(s))

	C.BannerNotify(t, m, s)

	return nil
}

// speechNotify triggers an NSSpeechSynthesizer notification.
func speechNotify(n notification) error {
	voice := os.Getenv(voiceEnv)
	if voice == "" {
		voice = "Alex"
	}
	text := fmt.Sprintf("%s %s", n.title, n.message)

	cmd := exec.Command("say", "-v", voice, text)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Speech: %s", err)
	}

	return nil
}
