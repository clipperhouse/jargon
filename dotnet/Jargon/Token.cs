namespace Jargon
{
    public readonly struct Token
    {
        private readonly string _String;
        public string String => _String;
        private readonly bool _Punct;
        public bool IsPunct => _Punct;
        private readonly bool _Space;
        public bool IsSpace => _Space;
        private readonly bool _Lemma;
        public bool IsLemma => _Lemma;

        internal Token(string value, bool punct, bool space, bool lemma)
        {
            _String = value;
            _Punct = punct;
            _Space = space;
            _Lemma = lemma;
        }

        public override string ToString() => String;
    }
}
