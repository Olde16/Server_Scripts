set -euo pipefail

# ==========================================================
#                    KONFIGURATION
# ==========================================================

# Remote (Quelle)
HOST="www.example.com"
USER="backup"
REMOTE="/var/www/html/data/"

# Lokales Backup-Ziel (ABSOLUTER PFAD!)
BASE="/BACKUP/www-backups"

# Backup-Strategie
FULL_EVERY_X_DAYS=14        # Vollbackup alle X Tage
DELETE_OLDER_THAN_DAYS=33   # Backups älter als X Tage löschen
# 14*3 = 32 > +1 sicherheit

# Logging
LOG_DIR="/BACKUP/Logs"

# ==========================================================
#                  VORBEREITUNG
# ==========================================================

mkdir -p "$BASE" "$LOG_DIR"

TS=$(date +"%Y-%m-%d_%H-%M-%S")
DEST="$BASE/$TS"
LAST="$BASE/latest"
LAST_FULL_FILE="$BASE/latest_full"
LOG="$LOG_DIR/backup_$TS.log"

# Alles loggen (stdout + stderr)
exec > >(tee -a "$LOG") 2>&1

echo "=============================================="
echo "Backup gestartet: $(date)"
echo "Ziel: $DEST"
echo "=============================================="

mkdir -p "$DEST"

# ==========================================================
#         ALTE BACKUPS SICHER LÖSCHEN (NAMENSBASIERT)
# ==========================================================

echo "Prüfe Backups älter als $DELETE_OLDER_THAN_DAYS Tage…"

NOW=$(date +%s)

for dir in "$BASE"/20*-*-*_*
do
    [[ -d "$dir" ]] || continue

    name=$(basename "$dir")
    date_part="${name:0:10}"

    if ! folder_ts=$(date -d "$date_part" +%s 2>/dev/null); then
        echo "Überspringe unbekanntes Verzeichnis: $dir"
        continue
    fi

    age_days=$(( (NOW - folder_ts) / 86400 ))

    if (( age_days > DELETE_OLDER_THAN_DAYS )); then
        echo "Lösche altes Backup ($age_days Tage): $dir"
        rm -rf "$dir"
    fi
done

# ==========================================================
#            PRÜFEN OB VOLLBACKUP NÖTIG
# ==========================================================

force_full=false

if [[ ! -f "$LAST_FULL_FILE" ]]; then
    echo "Kein vorheriges Vollbackup gefunden."
    force_full=true
else
    last_full_date=$(cat "$LAST_FULL_FILE")
    days_since_full=$(( ( $(date +%s) - $(date -d "$last_full_date" +%s) ) / 86400 ))

    echo "Letztes Vollbackup vor $days_since_full Tagen."

    if (( days_since_full >= FULL_EVERY_X_DAYS )); then
        echo "Vollbackup fällig."
        force_full=true
    fi
fi

# ==========================================================
#                    BACKUP
# ==========================================================

I_RSYNC_OPTS="-av --delete"
F_RSYNC_OPTS="-av --delete --checksum"

if [[ "$force_full" == true ]]; then
    echo "Starte VOLLBACKUP..."
    rsync $F_RSYNC_OPTS \
        "${USER}@${HOST}:${REMOTE}" \
        "$DEST"

    date +"%Y-%m-%d" > "$LAST_FULL_FILE"
else
    echo "Starte INKREMENTELLES BACKUP..."
    rsync $I_RSYNC_OPTS \
        --link-dest="$LAST" \
        "${USER}@${HOST}:${REMOTE}" \
        "$DEST"
fi

# ==========================================================
#               latest SYMLINK AKTUALISIEREN
# ==========================================================

ln -sfn "$DEST" "$LAST"

echo "=============================================="
echo "Backup abgeschlossen: $DEST"
echo "Logfile: $LOG"
echo "=============================================="
