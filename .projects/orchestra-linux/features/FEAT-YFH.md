---
created_at: "2026-02-28T02:55:39Z"
description: 'Set up gettext internationalization. po/POTFILES: list all .vala source files containing _() strings. po/LINGUAS: supported languages (en, fr, de, es, pt, ja, zh). Meson i18n module: i18n.gettext(''orchestra-desktop''). Mark all user-visible strings with _("text"). Extract strings: meson compile orchestra-desktop-pot. Merge: meson compile orchestra-desktop-update-po. Compile .po to .mo in build. Install to /usr/share/locale/LANG/LC_MESSAGES/orchestra-desktop.mo. Weblate integration config for community translations. GNOME Translation Project compatibility.'
id: FEAT-YFH
priority: P2
project_id: orchestra-linux
status: backlog
title: Gettext i18n / translations scaffolding
updated_at: "2026-02-28T02:55:39Z"
version: 0
---

# Gettext i18n / translations scaffolding

Set up gettext internationalization. po/POTFILES: list all .vala source files containing _() strings. po/LINGUAS: supported languages (en, fr, de, es, pt, ja, zh). Meson i18n module: i18n.gettext('orchestra-desktop'). Mark all user-visible strings with _("text"). Extract strings: meson compile orchestra-desktop-pot. Merge: meson compile orchestra-desktop-update-po. Compile .po to .mo in build. Install to /usr/share/locale/LANG/LC_MESSAGES/orchestra-desktop.mo. Weblate integration config for community translations. GNOME Translation Project compatibility.
