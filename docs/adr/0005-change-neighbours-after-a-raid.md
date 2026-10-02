# Bend the knee to change neighbours after a raid

With [Ring-order turns](0011-round-robin-turns.md), a successfully attacked
Target may Bend the Knee in the step immediately after Maintenance and before
Rumour. Their turn ends, their Ring position updates, and their original Target
begins their turn. All remaining phases are skipped, including every End Phase
step; no additional future turn is forfeited. This supersedes
the original choice and shared-resolution reordering timing described below;
the following text preserves the earlier rationale and batch-reordering model.

Attack impacts are already applied when the Target receives the attack report
at the start of their next turn; that report identifies what attacked them and
the outcome, and they are then offered the choice to bend the knee at no
additional Integrity cost, granting the Aggressor one Domination point and
changing Ring positions only after that next turn's committed attacks resolve.
Kneeling provides escape at the cost of a rival's victory progress, without
guaranteeing a better replacement Aggressor; there is no separate Integrity
payment option for escape.
Each connected chain of kneelings reverses its affected Ring order, and the
whole Ring reverses direction when everyone kneels; each submission grants
its point to the Aggressor receiving that submission.

## Example

In A → B → C → D, B kneeling to A and C kneeling to B produces C → B → A → D.

## Open questions

- How a previously chosen kneeling is handled if the intended Aggressor is
  eliminated before Ring reordering, or the resolution leaves only two survivors.
