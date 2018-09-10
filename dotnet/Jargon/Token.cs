using System;

namespace Jargon
{
    public readonly struct Token
    {
        [Flags]
        private enum State: byte
        {
            None = 0,
            Punct = 1 << 0,
            Space = 1 << 1,
            Lemma = 1 << 2
        }

        private readonly string _String;
        public string String => _String;

        private readonly State _State;

        public bool IsPunct => _State.HasFlag(State.Punct);
        public bool IsSpace => _State.HasFlag(State.Space);
        public bool IsLemma => _State.HasFlag(State.Lemma);

        internal Token(string value, bool punct, bool space, bool lemma)
        {
            _String = value;

            _State =
                (punct ? State.Punct : State.None) |
                (space ? State.Space : State.None) |
                (lemma ? State.Lemma : State.None);
        }

        public override string ToString() => String;
    }
}
